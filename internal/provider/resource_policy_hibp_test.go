package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyHaveIBeenPwned(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyHaveIBeenPwned(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_hibp.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyHaveIBeenPwned(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_hibp.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyHaveIBeenPwned(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_hibp" "name" {
  name              = "%[1]s"
}
`, name)
}
