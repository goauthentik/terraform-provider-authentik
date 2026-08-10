package sdkprovider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_users.admins", "users.0.username", "akadmin"),
					resource.TestCheckResourceAttr("data.authentik_users.admins", "users.0.is_superuser", "true"),
				),
			},
		},
	})
}

const testAccDataSourceUserSimple = `
data "authentik_users" "admins" {
  is_superuser = true
}
`
