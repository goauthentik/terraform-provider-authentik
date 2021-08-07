package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceScopeMapping(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceScopeMappingSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_scope_mapping.test", "name", "authentik default OAuth Mapping: Proxy outpost"),
					resource.TestCheckResourceAttr("data.authentik_scope_mapping.test", "managed", "goauthentik.io/providers/proxy/scope-proxy"),
				),
			},
		},
	})
}

const testAccDataSourceScopeMappingSimple = `
data "authentik_scope_mapping" "test" {
  name    = "authentik default OAuth Mapping: Proxy outpost"
  managed = "goauthentik.io/providers/proxy/scope-proxy"
}
`
