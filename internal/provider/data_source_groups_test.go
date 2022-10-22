package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroups(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupsSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.is_superuser", "true"),
				),
			},
		},
	})
}

const testAccDataSourceGroupsSimple = `
data "authentik_groups" "admins" {
  name = "authentik Admins"
}
`
