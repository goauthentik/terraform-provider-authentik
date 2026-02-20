package brand_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestAccDataSourceBrand(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
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
