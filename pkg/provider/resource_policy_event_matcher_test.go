package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyEventMatcher(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyEventMatcher(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_event_matcher.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyEventMatcher(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_event_matcher.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyEventMatcher(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_event_matcher" "name" {
  name              = "%[1]s"
  action = "login"
  app = "authentik.flows"
  client_ip = "1.2.3.4"
}
`, name)
}
