package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccDataSourceUsers(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
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
