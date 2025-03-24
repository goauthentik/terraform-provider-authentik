package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEventRule(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEventRule(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_event_rule.transport", "name", rName),
				),
			},
			{
				Config: testAccResourceEventRule(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_event_rule.transport", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceEventRule(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
}
resource "authentik_group" "group" {
  name = "%[1]s"
  users = [authentik_user.name.id]
  is_superuser = true
}
resource "authentik_event_transport" "transport" {
  name        = "%[1]s"
  mode        = "webhook_slack"
  send_once   = true
  webhook_url = "https://discord.com/...."
}
resource "authentik_event_rule" "transport" {
  name = "%[1]s"
  group = authentik_group.group.id
  transports = [authentik_event_transport.transport.id]
}
`, name)
}
