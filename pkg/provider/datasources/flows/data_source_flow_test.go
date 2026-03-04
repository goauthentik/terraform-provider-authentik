package flows_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestAccDataSourceFlow(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
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
