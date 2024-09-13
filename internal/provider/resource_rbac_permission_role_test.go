package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRBACRoleObjectPermission(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRBACRoleObjectPermissionScoped(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_permission_role.name", "permission", "authentik_core.view_application"),
				),
			},
			{
				Config: testAccRBACRoleObjectPermissionGlobal(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_permission_role.name", "permission", "authentik_core.add_application"),
				),
			},
		},
	})
}

func testAccRBACRoleObjectPermissionScoped(name string) string {
	return fmt.Sprintf(`
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}

resource "authentik_application" "name" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_rbac_permission_role" "name" {
  role = authentik_rbac_role.role.id
  model = "authentik_core.application"
  permission = "authentik_core.view_application"
  object_id = authentik_application.name.uuid
}
`, name)
}

func testAccRBACRoleObjectPermissionGlobal(name string) string {
	return fmt.Sprintf(`
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}

resource "authentik_rbac_permission_role" "name" {
  role = authentik_rbac_role.role.id
  permission = "authentik_core.add_application"
}
`, name)
}
