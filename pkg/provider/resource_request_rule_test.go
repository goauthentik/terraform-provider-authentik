package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRequestRule(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRequestRule(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_request_rule.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_request_rule.name", "min_reviewers", "1"),
				),
			},
			{
				Config: testAccResourceRequestRule(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_request_rule.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_request_rule.name", "min_reviewers", "2"),
				),
			},
		},
	})
}

func testAccResourceRequestRule(name string, minReviewers int) string {
	return fmt.Sprintf(`
resource "authentik_request_rule" "name" {
  name                        = "%[1]s"
  policy_engine_mode          = "any"
  notification_mode           = "all"
  min_reviewers               = %[2]d
  min_reviewers_is_per_group  = false
}
`, name, minReviewers)
}
