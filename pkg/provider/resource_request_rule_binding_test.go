package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRequestRuleBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRequestRuleBinding(rName, "hours=1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_request_rule_binding.binding", "expiry_pending", "hours=1"),
				),
			},
			{
				Config: testAccResourceRequestRuleBinding(rName, "hours=2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_request_rule_binding.binding", "expiry_pending", "hours=2"),
				),
			},
		},
	})
}

func testAccResourceRequestRuleBinding(name string, expiryPending string) string {
	return fmt.Sprintf(`
resource "authentik_request_rule" "name" {
  name = "%[1]s"
}
resource "authentik_application" "name" {
  name = "%[1]s"
  slug = "%[1]s"
}
resource "authentik_request_rule_binding" "binding" {
  rule           = authentik_request_rule.name.id
  target         = authentik_application.name.uuid
  expiry_pending = "%[2]s"
}
`, name, expiryPending)
}
