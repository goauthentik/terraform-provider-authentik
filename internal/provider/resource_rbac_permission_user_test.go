package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRBACUserObjectPermission(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRBACUserObjectPermissionScoped(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_permission_user.name", "permission", "authentik_core.view_application"),
				),
			},
			{
				Config: testAccRBACUserObjectPermissionGlobal(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_rbac_permission_user.name", "permission", "authentik_core.add_application"),
				),
			},
		},
	})
}

func testAccRBACUserObjectPermissionScoped(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "user" {
  username = "%[1]s"
}

resource "authentik_application" "name" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_rbac_permission_user" "name" {
  user = authentik_user.user.id
  model = "authentik_core.application"
  permission = "authentik_core.view_application"
  object_id = authentik_application.name.uuid
}
`, name)
}

func testAccRBACUserObjectPermissionGlobal(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "user" {
  username = "%[1]s"
}

resource "authentik_rbac_permission_user" "name" {
  user = authentik_user.user.id
  permission = "authentik_core.add_application"
}
`, name)
}
