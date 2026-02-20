package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccDataSourceBrand(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBrandSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_brand.authentik-default", "domain", "authentik-default"),
					resource.TestCheckResourceAttr("data.authentik_brand.authentik-default", "branding_title", "authentik"),
				),
			},
		},
	})
}

const testAccDataSourceBrandSimple = `
data "authentik_brand" "authentik-default" {
  domain = "authentik-default"
}
`
