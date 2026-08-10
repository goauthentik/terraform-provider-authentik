package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyEventMatcher(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
