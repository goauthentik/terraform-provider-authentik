package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyPassword(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyPassword(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_password.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyPassword(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_password.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyPassword(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_password" "name" {
  name              = "%[1]s"
  error_message = "foo"
}
`, name)
}
