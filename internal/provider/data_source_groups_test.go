package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroups(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupsSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.name", rName),
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins", "groups.0.users_obj.#", "1"),
				),
			},
			{
				Config: testAccDataSourceGroupsSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.name", rName),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_no_users_obj", "groups.0.users_obj.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGroupsSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.name", rName),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_groups.admins_include_users_obj", "groups.0.users_obj.#", "1"),
				),
			},
		},
	})
}

func testAccDataSourceGroupsSimple(name string) string {
	return fmt.Sprintf(`resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
}
resource "authentik_group" "group" {
  name = "%[1]s"
  users = [authentik_user.name.id]
  is_superuser = true
}

data "authentik_groups" "admins" {
  name = authentik_group.group.name
}

data "authentik_groups" "admins_no_users_obj" {
  name          = authentik_group.group.name
  include_users = false
}

data "authentik_groups" "admins_include_users_obj" {
  name          = authentik_group.group.name
  include_users = true
}`, name)
}
