package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyDummy(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyDummy(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyDummy(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyDummy(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_dummy" "name" {
  name              = "%[1]s"
}
`, name)
}
