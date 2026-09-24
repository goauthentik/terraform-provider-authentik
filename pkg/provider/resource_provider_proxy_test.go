package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccResourceProviderProxy(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderProxy(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_proxy.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_proxy.name", "external_host", "http://"+rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderProxy(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_proxy.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_proxy.name", "external_host", "http://"+rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderProxy(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_property_mapping_provider_scope" "test" {
  managed_list = [
    "goauthentik.io/providers/oauth2/scope-openid",
    "goauthentik.io/providers/oauth2/scope-email",
    "goauthentik.io/providers/oauth2/scope-profile",
    "goauthentik.io/providers/oauth2/scope-entitlements",
    "goauthentik.io/providers/proxy/scope-proxy",
  ]
}

resource "authentik_provider_proxy" "name" {
  name      = "%[1]s"
  internal_host = "http://foo.bar.baz"
  external_host = "http://%[1]s"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  property_mappings = data.authentik_property_mapping_provider_scope.test.ids
  skip_path_regex    = <<EOF
^/$
^/status
^/assets/
^/assets
^/icon.svg
^/api/.*
^/upload/.*
^/metrics
EOF
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_proxy.name.id
}
`, name, appName)
}

func TestResourceProviderProxyBooleanRequests(t *testing.T) {
	for _, tc := range []struct {
		name            string
		values          map[string]any
		sslValidation   bool
		basicAuth       bool
		interceptHeader bool
	}{
		{name: "defaults", sslValidation: true, interceptHeader: true},
		{
			name: "explicit false",
			values: map[string]any{
				"internal_host_ssl_validation": false,
				"basic_auth_enabled":           false,
				"intercept_header_auth":        false,
			},
		},
		{
			name: "explicit true",
			values: map[string]any{
				"internal_host_ssl_validation": true,
				"basic_auth_enabled":           true,
				"intercept_header_auth":        true,
			},
			sslValidation: true, basicAuth: true, interceptHeader: true,
		},
		{
			name: "null uses defaults",
			values: map[string]any{
				"internal_host_ssl_validation": nil,
				"basic_auth_enabled":           nil,
				"intercept_header_auth":        nil,
			},
			sslValidation: true, interceptHeader: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceProviderProxy().Schema, tc.values)
			body, err := json.Marshal(resourceProviderProxySchemaToProvider(d))
			require.NoError(t, err)

			var request map[string]any
			require.NoError(t, json.Unmarshal(body, &request))
			assert.Equal(t, tc.sslValidation, request["internal_host_ssl_validation"])
			assert.Equal(t, tc.basicAuth, request["basic_auth_enabled"])
			assert.Equal(t, tc.interceptHeader, request["intercept_header_auth"])
		})
	}
}

func TestAccResourceProviderProxyBooleans(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	steps := []resource.TestStep{}
	for _, tc := range []struct {
		attributes      string
		sslValidation   string
		basicAuth       string
		interceptHeader string
	}{
		{
			attributes: `internal_host_ssl_validation = false
  basic_auth_enabled = false
  intercept_header_auth = false`,
			sslValidation: "false", basicAuth: "false", interceptHeader: "false",
		},
		{
			attributes: `internal_host_ssl_validation = true
  basic_auth_enabled = true
  intercept_header_auth = true`,
			sslValidation: "true", basicAuth: "true", interceptHeader: "true",
		},
		{
			attributes: `internal_host_ssl_validation = false
  basic_auth_enabled = false
  intercept_header_auth = false`,
			sslValidation: "false", basicAuth: "false", interceptHeader: "false",
		},
		{sslValidation: "true", basicAuth: "false", interceptHeader: "true"},
	} {
		config := fmt.Sprintf(`
data "authentik_flow" "authorization" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "invalidation" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_proxy" "test" {
  name = %q
  internal_host = "https://backend.example.com"
  external_host = "https://proxy.example.com"
  basic_auth_username_attribute = "proxy_username"
  basic_auth_password_attribute = "proxy_password"
  authorization_flow = data.authentik_flow.authorization.id
  invalidation_flow = data.authentik_flow.invalidation.id
  %s
}
`, rName, tc.attributes)
		steps = append(steps, resource.TestStep{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("authentik_provider_proxy.test", "internal_host_ssl_validation", tc.sslValidation),
				resource.TestCheckResourceAttr("authentik_provider_proxy.test", "basic_auth_enabled", tc.basicAuth),
				resource.TestCheckResourceAttr("authentik_provider_proxy.test", "intercept_header_auth", tc.interceptHeader),
			),
		})
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps:             steps,
	})
}
