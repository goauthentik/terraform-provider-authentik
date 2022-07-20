package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTenant(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTenantSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_tenant.authentik-default", "domain", "authentik-default"),
					resource.TestCheckResourceAttr("data.authentik_tenant.authentik-default", "branding_title", "authentik"),
				),
			},
		},
	})
}

const testAccDataSourceTenantSimple = `
data "authentik_stage" "authentik-default" {
  domain = "authentik-default"
}
`
