package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserLogin(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogin(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "session_duration", "minutes=1"),
				),
			},
			{
				Config: testAccResourceStageUserLogin(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "session_duration", "minutes=1"),
				),
			},
		},
	})
}

func testAccResourceStageUserLogin(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_login" "name" {
  name              = "%[1]s"
  session_duration = "minutes=1"
}
`, name)
}
