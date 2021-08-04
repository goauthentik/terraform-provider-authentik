package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageInvitation(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageInvitation(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_invitation.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageInvitation(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_invitation.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageInvitation(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_invitation" "name" {
  name              = "%[1]s"
}
`, name)
}
