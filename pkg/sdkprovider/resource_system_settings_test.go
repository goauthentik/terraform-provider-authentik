package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSystemSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSystemSettings("true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_system_settings.name", "default_user_change_username", "true"),
				),
			},
			{
				Config: testAccResourceSystemSettings("false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_system_settings.name", "default_user_change_username", "false"),
				),
			},
		},
	})
}

func testAccResourceSystemSettings(value string) string {
	return fmt.Sprintf(`
resource "authentik_system_settings" "name" {
  footer_links = [
	{
		name = "test"
		href = "https://google.com"
	}
  ]
  default_user_change_username = %[1]s
}
`, value)
}
