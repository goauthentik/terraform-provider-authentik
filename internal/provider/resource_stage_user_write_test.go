package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserWrite(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserWrite(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "create_users_as_inactive", "false"),
				),
			},
			{
				Config: testAccResourceStageUserWrite(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "create_users_as_inactive", "false"),
				),
			},
		},
	})
}

func testAccResourceStageUserWrite(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_write" "name" {
  name              = "%[1]s"
  create_users_as_inactive = false
}
`, name)
}
