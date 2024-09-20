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
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.users_obj.#", "1"),
				),
			},
			{
				Config: testAccDataSourceGroupsSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.users_obj.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGroupsSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.users_obj.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceGroupsSimple = `
data "authentik_groups" "admins" {
  name = "authentik Admins"
}

data "authentik_groups" "admins_no_users_obj" {
  name          = "authentik Admins"
  include_users = false
}

data "authentik_groups" "admins_include_users_obj" {
  name          = "authentik Admins"
  include_users = true
}
`
