package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUserOffboarding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	scheduledAt := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserOffboarding(rName, scheduledAt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user_offboarding.offboarding", "action", "deactivate"),
					resource.TestCheckResourceAttrSet("authentik_user_offboarding.offboarding", "status"),
				),
			},
		},
	})
}

func testAccResourceUserOffboarding(name string, scheduledAt string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name     = "%[1]s"
}

resource "authentik_user_offboarding" "offboarding" {
  user            = authentik_user.name.id
  scheduled_at    = "%[2]s"
  action          = "deactivate"
  revoke_sessions = true
  revoke_tokens   = true
}
`, name, scheduledAt)
}
