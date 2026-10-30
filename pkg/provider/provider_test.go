package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	api "goauthentik.io/api/v3"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"authentik": func() (*schema.Provider, error) {
		return Provider("test", false), nil
	},
}

var providerTestFactories = map[string]func() (*schema.Provider, error){
	"authentik": func() (*schema.Provider, error) {
		return Provider("test", true), nil
	},
}

func TestProvider(t *testing.T) {
	p := Provider("testing", false)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %[1]s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	testEnvIsSet("AUTHENTIK_URL", t)
	testEnvIsSet("AUTHENTIK_TOKEN", t)
}

// testAccPreCheckEnterprise skips the test when the authentik instance has no valid enterprise license
// (e.g. CI runs for PRs from forks, which don't get access to the license secret)
func testAccPreCheckEnterprise(t *testing.T) {
	testAccPreCheck(t)
	p := Provider("test", false)
	ac, diags := p.ConfigureContextFunc(t.Context(), schema.TestResourceDataRaw(t, p.Schema, map[string]any{}))
	if diags.HasError() {
		t.Fatalf("failed to configure provider: %v", diags)
	}
	summary, _, err := ac.(*APIClient).client.EnterpriseAPI.EnterpriseLicenseSummaryRetrieve(t.Context()).Execute()
	if err != nil {
		t.Fatalf("failed to get license summary: %v", err)
	}
	if summary.Status != api.LICENSESUMMARYSTATUSENUM_VALID && summary.Status != api.LICENSESUMMARYSTATUSENUM_EXPIRY_SOON {
		t.Skipf("enterprise license not available (status: %s)", summary.Status)
	}
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
				"token":    "",
				"insecure": false,
			}))
			assert.Nil(t, diag)
			ac := _ac.(*APIClient)
			assert.Equal(t, tc.expectedURL, ac.client.GetConfig().Servers[0].URL, "Server URL should be constructed correctly")
		})
	}
}
