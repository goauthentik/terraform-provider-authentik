package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserSimplePk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "username", "datasource-test"),
					resource.TestCheckResourceAttr("data.authentik_user.akadmin", "is_superuser", "false"),
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
resource "authentik_user" "test" {
  username = "datasource-test"
}

data "authentik_user" "akadmin" {
  pk = authentik_user.test.id
}
`
