package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

const testCertificateKeyPairPEM = "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"

func testCertificateKeyPairServer(t *testing.T, keyData string) *APIClient {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/crypto/certificatekeypairs/kp/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"pk":"kp","name":"sp-verification","fingerprint_sha256":"a","fingerprint_sha1":"b","cert_expiry":"2030-01-01T00:00:00Z","cert_subject":"CN=test","private_key_available":` + map[bool]string{true: "true", false: "false"}[keyData != ""] + `,"key_type":null,"certificate_download_url":"/c","private_key_download_url":"/k","managed":null}`))
	})
	mux.HandleFunc("/api/v3/crypto/certificatekeypairs/kp/view_certificate/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":"-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"}`))
	})
	mux.HandleFunc("/api/v3/crypto/certificatekeypairs/kp/view_private_key/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":"` + keyData + `"}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cfg := api.NewConfiguration()
	cfg.Servers = api.ServerConfigurations{{URL: srv.URL + "/api/v3"}}
	cfg.HTTPClient = srv.Client()
	return &APIClient{client: api.NewAPIClient(cfg)}
}

// A certificate imported without a private key comes back from authentik with
// an empty `data` on view_private_key. That must not end up in state as "\n",
// otherwise every plan proposes `key_data = (sensitive value) -> null`.
func TestResourceCertificateKeyPairReadWithoutKeyNoDiff(t *testing.T) {
	client := testCertificateKeyPairServer(t, "")
	r := resourceCertificateKeyPair()

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		"name":             "sp-verification",
		"certificate_data": testCertificateKeyPairPEM + "\n",
	})
	d.SetId("kp")

	diags := resourceCertificateKeyPairRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "", d.Get("key_data"))

	state := d.State()
	cfg := terraform.NewResourceConfigRaw(map[string]any{
		"name":             "sp-verification",
		"certificate_data": testCertificateKeyPairPEM + "\n",
	})
	diff, err := r.Diff(context.Background(), state, cfg, nil)
	require.NoError(t, err)
	if diff != nil {
		_, ok := diff.Attributes["key_data"]
		assert.False(t, ok, "unexpected diff on key_data: %v", diff.Attributes["key_data"])
	}
}

func TestResourceCertificateKeyPairReadWithKey(t *testing.T) {
	client := testCertificateKeyPairServer(t, "-----BEGIN PRIVATE KEY-----\\nMIIE\\n-----END PRIVATE KEY-----")
	r := resourceCertificateKeyPair()

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
		"name":             "sp-verification",
		"certificate_data": testCertificateKeyPairPEM + "\n",
	})
	d.SetId("kp")

	diags := resourceCertificateKeyPairRead(context.Background(), d, client)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, "-----BEGIN PRIVATE KEY-----\nMIIE\n-----END PRIVATE KEY-----\n", d.Get("key_data"))
}
