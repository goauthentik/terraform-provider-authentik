package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRBACRole(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRBACRole(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_role.group", "name", rName),
				),
			},
			{
				Config: testAccRBACRole(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_role.group", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccRBACRole(name string) string {
	return fmt.Sprintf(`
resource "authentik_rbac_role" "group" {
  name = "%[1]s"
}
`, name)
}
