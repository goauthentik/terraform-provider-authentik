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
				Config: testAccResourcePolicyBindingPolicy(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyBindingPolicy(rName+"test", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_dummy.name", "name", rName+"test"),
				),
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyBindingGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_group.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyBindingGroup(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_group.name", "name", rName+"test"),
				),
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyBindingUser(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName),
				),
			},
			{
				Config: testAccResourcePolicyBindingUser(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_user.name", "username", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyBindingPolicy(name string, order int) string {
	return fmt.Sprintf(`
resource "authentik_policy_dummy" "name" {
  name              = "%[1]s"
}
resource "authentik_application" "name" {
  name              = "%[1]s-policy"
  slug              = "%[1]s-policy"
}
resource "authentik_policy_binding" "binding" {
  target = authentik_application.name.uuid
  policy = authentik_policy_dummy.name.id
  order = %[2]d
}
`, name, order)
}

func testAccResourcePolicyBindingGroup(name string) string {
	return fmt.Sprintf(`
resource "authentik_group" "name" {
  name              = "%[1]s"
}
resource "authentik_application" "name" {
  name              = "%[1]s-group"
  slug              = "%[1]s-group"
}
resource "authentik_policy_binding" "binding" {
  target = authentik_application.name.uuid
  group = authentik_group.name.id
  order = 0
}
`, name)
}

func testAccResourcePolicyBindingUser(name string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
  username              = "%[1]s"
}
resource "authentik_application" "name" {
  name              = "%[1]s-user"
  slug              = "%[1]s-user"
}
resource "authentik_policy_binding" "binding" {
  target = authentik_application.name.uuid
  user = authentik_user.name.id
  order = 0
}
`, name)
}
