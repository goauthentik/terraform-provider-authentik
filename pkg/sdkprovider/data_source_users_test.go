package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUsers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
