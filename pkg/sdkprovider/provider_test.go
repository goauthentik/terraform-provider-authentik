package sdkprovider

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	testinginterface "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
	fwprovider "goauthentik.io/terraform-provider-authentik/pkg/provider"
)

// muxProviderServer builds the same muxed protocol 6 server main.go serves, so
// acceptance tests exercise the exact composition users get rather than the bare
// SDKv2 provider in isolation.
func muxProviderServer(version string, testing bool) (tfprotov6.ProviderServer, error) {
	upgradedSDKServer, err := tf5to6server.UpgradeServer(context.Background(), Provider(version, testing).GRPCProvider)
	if err != nil {
		return nil, err
	}

	return tf6muxserver.NewMuxServer(context.Background(),
		func() tfprotov6.ProviderServer { return upgradedSDKServer },
		providerserver.NewProtocol6(fwprovider.New(version, testing)),
	)
}

// ProtoV6ProviderFactories are used to instantiate the muxed provider during
// acceptance testing. The factory function will be invoked for every Terraform CLI
// command executed to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"authentik": func() (tfprotov6.ProviderServer, error) {
		return muxProviderServer("test", false)
	},
}

var providerTestFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"authentik": func() (tfprotov6.ProviderServer, error) {
		return muxProviderServer("test", true)
	},
}

func TestProvider(t *testing.T) {
	p := Provider("testing", false)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %[1]s", err)
	}
}

// schemaCmpOptions mirrors tf6muxserver's own (unexported) schema comparison options:
// SchemaAttribute/SchemaNestedBlock slices may differ in order without being a real
// mismatch, and MinItems/MaxItems on nested blocks are excluded from mux's own equality
// check. See github.com/hashicorp/terraform-plugin-mux/tf6muxserver/schema_equality.go.
var schemaCmpOptions = []cmp.Option{
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaAttribute) bool {
		return i.Name < j.Name
	}),
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaNestedBlock) bool {
		return i.TypeName < j.TypeName
	}),
	cmpopts.IgnoreFields(tfprotov6.SchemaNestedBlock{}, "MinItems", "MaxItems"),
}

// TestProviderSchemaMatchesSDKv2 guards the mux-compatibility work in 0a: if the SDKv2
// provider block schema and the framework provider block schema ever diverge,
// tf6muxserver.NewMuxServer's own comparison would reject the pairing at startup with
// "Invalid Provider Server Combination" - this test catches that at `go test` time
// instead of `terraform plan` time.
func TestProviderSchemaMatchesSDKv2(t *testing.T) {
	ctx := context.Background()

	sdkServer, err := tf5to6server.UpgradeServer(ctx, Provider("testing", false).GRPCProvider)
	require.NoError(t, err)
	sdkSchema, err := sdkServer.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, sdkSchema.Diagnostics)

	fwServer := providerserver.NewProtocol6(fwprovider.New("testing", false))()
	fwSchema, err := fwServer.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, fwSchema.Diagnostics)

	if diff := cmp.Diff(sdkSchema.Provider, fwSchema.Provider, schemaCmpOptions...); diff != "" {
		t.Errorf("provider block schema mismatch between pkg/sdkprovider and pkg/provider (-sdk +framework):\n%s", diff)
	}
}

func testAccPreCheck(t *testing.T) {
	testEnvIsSet("AUTHENTIK_URL", t)
	testEnvIsSet("AUTHENTIK_TOKEN", t)
}

// testAccAPIClientFromEnv builds a live API client from AUTHENTIK_URL/AUTHENTIK_TOKEN,
// for use in CheckDestroy functions, which only receive a *terraform.State and have
// no *testing.T to plumb through.
func testAccAPIClientFromEnv() *api.APIClient {
	rt := &testinginterface.RuntimeT{}
	p := Provider("test", false)
	raw, diags := p.ConfigureContextFunc(context.Background(), schema.TestResourceDataRaw(rt, p.Schema, map[string]any{}))
	if diags.HasError() {
		rt.Fatalf("failed to configure provider: %v", diags)
	}
	return raw.(*APIClient).client
}

func testEnvIsSet(k string, t *testing.T) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("%[1]s must be set for acceptance tests", k)
	}
}

func TestProviderConfigure_PathBasedURL(t *testing.T) {
	testCases := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "Root path with trailing slash",
			inputURL:    "https://api.example.com/",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Root path without trailing slash",
			inputURL:    "https://api.example.com",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Single segment path with trailing slash",
			inputURL:    "https://api.example.com/sso/",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Single segment path without trailing slash",
			inputURL:    "https://api.example.com/sso",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Multi-segment path with trailing slash",
			inputURL:    "https://api.example.com/auth/sso/",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "Multi-segment path without trailing slash",
			inputURL:    "https://api.example.com/auth/sso",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "HTTP scheme with path",
			inputURL:    "http://localhost:9000/sso/",
			expectedURL: "http://localhost:9000/sso/api/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Provider("testing", true)

			_ac, diag := p.ConfigureContextFunc(t.Context(), schema.TestResourceDataRaw(t, p.Schema, map[string]any{
				"url":      tc.inputURL,
				"token":    "test-token",
				"insecure": false,
			}))
			assert.Nil(t, diag)
			ac := _ac.(*APIClient)
			assert.Equal(t, tc.expectedURL, ac.client.GetConfig().Servers[0].URL, "Server URL should be constructed correctly")
		})
	}
}
