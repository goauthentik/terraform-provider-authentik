package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicy(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicySimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy.default-authentication-flow-password-stage", "name", "default-authentication-flow-password-stage"),
				),
			},
		},
	})
}

const testAccDataSourcePolicySimple = `
data "authentik_policy" "default-authentication-flow-password-stage" {
  name = "default-authentication-flow-password-stage"
}
`
