package sources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestAccDataSourceSource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
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
