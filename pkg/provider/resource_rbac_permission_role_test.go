package provider

import (
	"fmt"
	"strings"
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

// TestAccRBACRoleObjectPermissionPagination tests that permissions on page 2+
// are correctly found when a role has more than 20 permissions assigned.
// This is a regression test for the pagination bug where only the first page
// of results was fetched.
func TestAccRBACRoleObjectPermissionPagination(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRBACRoleObjectPermissionManyGlobal(rName),
				Check: resource.ComposeTestCheckFunc(
					// Check a permission likely on page 1
					resource.TestCheckResourceAttr("authentik_rbac_permission_role.perms-authentik_core_add_user", "permission", "authentik_core.add_user"),
					// Check a permission likely on page 2 (beyond default page size of 20)
					resource.TestCheckResourceAttr("authentik_rbac_permission_role.perms-authentik_providers_oauth2_view_oauth2provider", "permission", "authentik_providers_oauth2.view_oauth2provider"),
				),
			},
		},
	})
}

func testAccRBACRoleObjectPermissionManyGlobal(name string) string {
	// 21 permissions to ensure we exceed the default page size of 20
	perms := []string{
		"authentik_core.add_user",
		"authentik_core.change_user",
		"authentik_core.delete_user",
		"authentik_core.view_user",
		"authentik_core.add_group",
		"authentik_core.change_group",
		"authentik_core.delete_group",
		"authentik_core.view_group",
		"authentik_core.add_application",
		"authentik_core.change_application",
		"authentik_core.delete_application",
		"authentik_core.view_application",
		"authentik_core.add_token",
		"authentik_core.change_token",
		"authentik_core.delete_token",
		"authentik_core.view_token",
		"authentik_flows.add_flow",
		"authentik_flows.change_flow",
		"authentik_flows.delete_flow",
		"authentik_flows.view_flow",
		"authentik_providers_oauth2.view_oauth2provider",
	}
	mf := fmt.Sprintf(`
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}

`, name)
	for _, perm := range perms {
		mf = mf + fmt.Sprintf(`
resource "authentik_rbac_permission_role" "perms-%s" {
  role       = authentik_rbac_role.role.id
  permission = "%s"
}
`, strings.ReplaceAll(perm, ".", "_"), perm)
	}
	return mf
}
