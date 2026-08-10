package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceBrand(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
