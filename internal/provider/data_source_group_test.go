package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_group.test", "name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_group.test", "is_superuser", "true"),
				),
			},
		},
	})
}

const testAccDataSourceGroupSimple = `
data "authentik_group" "test" {
  name = "authentik Admins"
}
`
