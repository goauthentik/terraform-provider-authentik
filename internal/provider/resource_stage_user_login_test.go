package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserLogin(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogin,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "name", "test app"),
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "session_duration", "minutes=1"),
				),
			},
		},
	})
}

const testAccResourceStageUserLogin = `
resource "authentik_stage_user_login" "name" {
  name              = "test app"
  session_duration = "minutes=1"
}
`
