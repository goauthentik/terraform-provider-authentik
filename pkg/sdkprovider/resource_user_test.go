package sdkprovider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceUser(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
				),
			},
			{
				Config: testAccResourceUser(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
				),
			},
			{
				ResourceName:            "authentik_user.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccResourceUserGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
				),
			},
			{
				Config: testAccResourceUserGroup(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
					resource.TestCheckResourceAttr("authentik_user.name", "groups.#", "1"),
				),
			},
			{
				Config: testAccResourceUserRole(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
				),
			},
			{
				Config: testAccResourceUserRole(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
					resource.TestCheckResourceAttr("authentik_user.name", "roles.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceUserAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserAttributes(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_user.name", "attributes", `{"foo":"bar"}`),
				),
			},
		},
	})
}

func testAccCheckUserDestroy(s *terraform.State) error {
	c := testAccAPIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_user" {
			continue
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}
		_, hr, err := c.CoreAPI.CoreUsersRetrieve(context.Background(), int32(id)).Execute()
		if err == nil {
			return fmt.Errorf("user %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceUser(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
  password = "%[1]s"
}
`, name)
}

func testAccResourceUserGroup(name string) string {
	return fmt.Sprintf(`
resource "authentik_group" "group" {
  name = "%[1]s"
  is_superuser = true
}
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
  groups = [authentik_group.group.id]
}
`, name)
}

func testAccResourceUserRole(name string) string {
	return fmt.Sprintf(`
resource "authentik_group" "group" {
  name = "%[1]s"
  is_superuser = true
}
resource "authentik_rbac_role" "role" {
  name = "%[1]s"
}
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
  roles = [authentik_rbac_role.role.id]
}
`, name)
}

func testAccResourceUserAttributes(name string, valid bool) string {
	attributes := "jsonencode({\"foo\"= \"bar\"})"
	if !valid {
		attributes = "\"{\""
	}
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username = "%[1]s"
  name = "%[1]s"
  attributes = %[2]s
}
`, name, attributes)
}
