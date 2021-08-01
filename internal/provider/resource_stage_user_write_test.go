package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserWrite(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserWrite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "name", "test app"),
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "create_users_as_inactive", "false"),
				),
			},
		},
	})
}

const testAccResourceStageUserWrite = `
resource "authentik_stage_user_write" "name" {
  name              = "test app"
  create_users_as_inactive = false
}
`
