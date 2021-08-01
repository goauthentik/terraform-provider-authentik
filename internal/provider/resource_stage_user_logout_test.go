package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserLogout(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_logout.name", "name", "test"),
				),
			},
		},
	})
}

const testAccResourceStageUserLogout = `
resource "authentik_stage_user_logout" "name" {
  name              = "test"
}
`
