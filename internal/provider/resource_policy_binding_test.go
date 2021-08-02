package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyBinding(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyBinding(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyBinding(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_dummy" "name" {
  name              = "%s"
}
resource "authentik_application" "name" {
  name              = "%[1]s"
  slug              = "%[1]s"
}
resource "authentik_policy_binding" "binding" {
  target = authentik_application.name.uuid
  policy = authentik_policy_dummy.name.id
  order = 0
}
`, name)
}
