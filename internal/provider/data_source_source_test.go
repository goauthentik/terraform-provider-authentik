package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
