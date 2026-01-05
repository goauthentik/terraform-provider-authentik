package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRBACInitialPermissions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRBACInitialPermissions(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_role.group", "name", rName),
				),
			},
			{
				Config: testAccRBACInitialPermissions(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_role.group", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccRBACInitialPermissions(name string) string {
	return fmt.Sprintf(`
resource "authentik_rbac_role" "group" {
  name = "%[1]s"
}
data "authentik_rbac_permission" "admin" {
  codename = "access_admin_interface"
}
resource "authentik_rbac_initial_permissions" "ip" {
  name = "%[1]s"
  role = authentik_rbac_role.group.id
  permissions = [
    data.authentik_rbac_permission.admin.id,
  ]
}
`, name)
}
