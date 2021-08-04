package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageUserLogout(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogout(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_logout.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageUserLogout(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_logout.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageUserLogout(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_logout" "name" {
  name              = "%[1]s"
}
`, name)
}
