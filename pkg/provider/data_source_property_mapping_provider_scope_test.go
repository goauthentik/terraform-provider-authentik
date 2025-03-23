package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePropertyMappingProviderScope(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropertyMappingProviderScopeSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scope.test", "name", "authentik default OAuth Mapping: Proxy outpost"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scope.test", "managed", "goauthentik.io/providers/proxy/scope-proxy"),
				),
			},
			{
				Config: testAccDataSourceScopePropertyMappingList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scope.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourcePropertyMappingProviderScopeSimple = `
data "authentik_property_mapping_provider_scope" "test" {
  name    = "authentik default OAuth Mapping: Proxy outpost"
  managed = "goauthentik.io/providers/proxy/scope-proxy"
}
`

const testAccDataSourceScopePropertyMappingList = `
data "authentik_property_mapping_provider_scope" "test" {
  managed_list = [
    "goauthentik.io/providers/oauth2/scope-openid",
    "goauthentik.io/providers/oauth2/scope-email"
  ]
}
`
