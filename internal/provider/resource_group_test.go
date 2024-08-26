package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_group.group", "name", rName),
				),
			},
			{
				Config: testAccResourceGroup(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
					resource.TestCheckResourceAttr("authentik_group.group", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceGroup(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
}
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}
resource "authentik_group" "group" {
  name = "%[1]s"
  users = [authentik_user.name.id]
  is_superuser = true
  roles = [authentik_rbac_role.role.id]
}
`, name)
}
