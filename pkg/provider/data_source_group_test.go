package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroupSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_group.test", "name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_group.test", "is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_group.test", "users_obj.#", "1"),
				),
			},
			{
				Config: testAccDataSourceGroupSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_group.test_no_users_obj", "name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_group.test_no_users_obj", "is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_group.test_no_users_obj", "users_obj.#", "0"),
				),
			},
			{
				Config: testAccDataSourceGroupSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_group.test_include_users_obj", "name", "authentik Admins"),
					resource.TestCheckResourceAttr("data.authentik_group.test_include_users_obj", "is_superuser", "true"),
					resource.TestCheckResourceAttr("data.authentik_group.test_include_users_obj", "users_obj.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceGroupSimple = `
data "authentik_group" "test" {
  name = "authentik Admins"
}

data "authentik_group" "test_no_users_obj" {
  name          = "authentik Admins"
  include_users = false
}

data "authentik_group" "test_include_users_obj" {
  name          = "authentik Admins"
  include_users = true
}
`
