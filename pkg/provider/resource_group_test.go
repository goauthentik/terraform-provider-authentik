package provider_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
)

func TestAccResourceGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckGroupDestroy,
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
			{
				ResourceName:      "authentik_group.group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGroupDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_group" {
			continue
		}
		_, hr, err := c.CoreAPI.CoreGroupsRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("group %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
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
