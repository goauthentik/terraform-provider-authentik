package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUser(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
		},
	})
}

func TestAccResourceUserAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserAttributes(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
					resource.TestCheckResourceAttr("authentik_user.name", "attributes", "{\"foo\":\"bar\"}"),
				),
			},
			{
				Config:      testAccResourceUserAttributes(rName, false),
				ExpectError: regexp.MustCompile("unexpected EOF"),
			},
		},
	})
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
