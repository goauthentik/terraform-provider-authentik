package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAgent(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAgent(rName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_agent.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_agent.name", "is_active", "false"),
					resource.TestCheckResourceAttrSet("authentik_agent.name", "token"),
					resource.TestCheckResourceAttrSet("authentik_agent.name", "uuid"),
				),
			},
			{
				Config: testAccResourceAgent(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_agent.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_agent.name", "is_active", "true"),
				),
			},
		},
	})
}

func testAccResourceAgent(name string, isActive string) string {
	return fmt.Sprintf(`
resource "authentik_user" "parent" {
  username = "%[1]s-parent"
  name     = "%[1]s-parent"
}

resource "authentik_agent" "name" {
  username        = "%[1]s"
  name            = "%[1]s"
  parent          = authentik_user.parent.id
  policy_behavior = "mirror"
  is_active       = %[2]s
}
`, name, isActive)
}
