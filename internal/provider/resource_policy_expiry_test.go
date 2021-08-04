package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyExpiry(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyExpiry(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_expiry.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyExpiry(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_expiry.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyExpiry(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_expiry" "name" {
  name              = "%[1]s"
  days = 3
}
`, name)
}
