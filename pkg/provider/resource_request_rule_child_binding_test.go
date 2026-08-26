package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRequestRuleChildBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRequestRuleChildBinding(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("authentik_request_rule_child_binding.child", "binding"),
					resource.TestCheckResourceAttrSet("authentik_request_rule_child_binding.child", "target"),
				),
			},
			{
				Config: testAccResourceRequestRuleChildBinding(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("authentik_request_rule_child_binding.child", "binding"),
					resource.TestCheckResourceAttrSet("authentik_request_rule_child_binding.child", "target"),
				),
			},
		},
	})
}

func testAccResourceRequestRuleChildBinding(name string) string {
	return fmt.Sprintf(`
resource "authentik_request_rule" "name" {
  name = "%[1]s"
}
resource "authentik_application" "parent" {
  name = "%[1]s-parent"
  slug = "%[1]s-parent"
}
resource "authentik_application" "child" {
  name = "%[1]s-child"
  slug = "%[1]s-child"
}
resource "authentik_request_rule_binding" "binding" {
  rule   = authentik_request_rule.name.id
  target = authentik_application.parent.uuid
}
resource "authentik_request_rule_child_binding" "child" {
  binding = authentik_request_rule_binding.binding.id
  target  = authentik_application.child.uuid
}
`, name)
}
