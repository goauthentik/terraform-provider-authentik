package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserSimplePk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "username", "akadmin"),
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "is_superuser", "true"),
				),
			},
			{
				Config: testAccDataSourceUserSimpleUsername,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "username", "akadmin"),
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "is_superuser", "true"),
				),
			},
		},
	})
}

const testAccDataSourceUserSimpleUsername = `
data "authentik_user" "akadmin" {
  username = "akadmin"
}
`

const testAccDataSourceUserSimplePk = `
data "authentik_user" "akadmin" {
  pk = 1
}
`
