package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSourceSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_source.inbuilt", "managed", "goauthentik.io/sources/inbuilt"),
					resource.TestCheckResourceAttrSet("data.authentik_source.inbuilt", "uuid"),
				),
			},
		},
	})
}

const testAccDataSourceSourceSimple = `
data "authentik_source" "inbuilt" {
  managed = "goauthentik.io/sources/inbuilt"
}
`
