package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSystemSettings(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
