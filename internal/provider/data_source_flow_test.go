package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFlow(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlowSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_flow.default-authorization-flow", "slug", "default-provider-authorization-implicit-consent"),
					resource.TestCheckResourceAttr("data.authentik_flow.default-authorization-flow", "designation", "authorization"),
					resource.TestCheckResourceAttr("data.authentik_flow.default-authorization-flow", "authentication", "require_authenticated"),
				),
			},
		},
	})
}

const testAccDataSourceFlowSimple = `
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
`
