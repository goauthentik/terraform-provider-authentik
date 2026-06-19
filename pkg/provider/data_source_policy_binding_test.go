package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyBindingIdConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_binding.id-test", "order", "10"),
				),
			},
			{
				Config: testAccDataSourcePolicyBindingGroupConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_binding.group-test1", "order", "10"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.group-test2-20", "order", "20"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.group-test2-30", "order", "30"),
				),
			},
			{
				Config: testAccDataSourcePolicyBindingPolicyConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_binding.policy-test1", "order", "10"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.policy-test2-20", "order", "20"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.policy-test2-30", "order", "30"),
				),
			},
			{
				Config: testAccDataSourcePolicyBindingUserConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_binding.user-test1", "order", "10"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.user-test2-20", "order", "20"),
					resource.TestCheckResourceAttr("data.authentik_policy_binding.user-test2-30", "order", "30"),
				),
			},
		},
	})
}

func TestAccDataSourcePolicyBinding_Errors(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePolicyBindingNoConfigSimple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourcePolicyBindingBrokenConfig1Simple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourcePolicyBindingBrokenConfig2Simple,
				ExpectError: regexp.MustCompile(`Neither id nor target and user/group/policy were provided`),
			},
			{
				Config:      testAccDataSourcePolicyBindingNotFound(appName, rName),
				ExpectError: regexp.MustCompile(`No matching policy bindings found`),
			},
			{
				Config:      testAccDataSourcePolicyBindingMultiple(appName, rName),
				ExpectError: regexp.MustCompile(`Multiple matching policy bindings found. Use order to select one.`),
			},
		},
	})
}

func testAccDataSourcePolicyBindingIdConfig(appName string, groupName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_group" "test" {
  name = "%[2]s"
}

resource "authentik_policy_binding" "group-test" {
  group  = authentik_group.test.id
  target = authentik_application.test.uuid
  order  = 10
}

data "authentik_policy_binding" "group-test" {
  group  = authentik_policy_binding.group-test.group
  target = authentik_policy_binding.group-test.target
}

data "authentik_policy_binding" "id-test" {
  id = data.authentik_policy_binding.group-test.id
}

`, appName, groupName)
}

func testAccDataSourcePolicyBindingGroupConfig(appName string, groupName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_group" "test1" {
  name = "%[2]s-1"
}

resource "authentik_group" "test2" {
  name = "%[2]s-2"
}

resource "authentik_policy_binding" "group-test1" {
  group  = authentik_group.test1.id
  target = authentik_application.test.uuid
  order  = 10
}

resource "authentik_policy_binding" "group-test2-20" {
  group  = authentik_group.test2.id
  target = authentik_application.test.uuid
  order  = 20
}

resource "authentik_policy_binding" "group-test2-30" {
  group  = authentik_group.test2.id
  target = authentik_application.test.uuid
  order  = 30
}

data "authentik_policy_binding" "group-test1" {
  group  = authentik_policy_binding.group-test1.group
  target = authentik_policy_binding.group-test1.target
}

data "authentik_policy_binding" "group-test2-20" {
  group  = authentik_policy_binding.group-test2-20.group
  target = authentik_policy_binding.group-test2-20.target
  order  = authentik_policy_binding.group-test2-20.order
}

data "authentik_policy_binding" "group-test2-30" {
  group  = authentik_policy_binding.group-test2-30.group
  target = authentik_policy_binding.group-test2-30.target
  order  = authentik_policy_binding.group-test2-30.order
}
`, appName, groupName)
}

func testAccDataSourcePolicyBindingPolicyConfig(appName string, policyName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_policy_expression" "test1" {
  name = "%[2]s-1"
  expression = "return True"
}

resource "authentik_policy_expression" "test2" {
  name = "%[2]s-2"
  expression = "return True"
}

resource "authentik_policy_binding" "policy-test1" {
  policy = authentik_policy_expression.test1.id
  target = authentik_application.test.uuid
  order  = 10
}

resource "authentik_policy_binding" "policy-test2-20" {
  policy = authentik_policy_expression.test2.id
  target = authentik_application.test.uuid
  order  = 20
}

resource "authentik_policy_binding" "policy-test2-30" {
  policy = authentik_policy_expression.test2.id
  target = authentik_application.test.uuid
  order  = 30
}

data "authentik_policy_binding" "policy-test1" {
  policy = authentik_policy_binding.policy-test1.policy
  target = authentik_policy_binding.policy-test1.target
  order  = authentik_policy_binding.policy-test1.order
}

data "authentik_policy_binding" "policy-test2-20" {
  policy = authentik_policy_binding.policy-test2-20.policy
  target = authentik_policy_binding.policy-test2-20.target
  order  = authentik_policy_binding.policy-test2-20.order
}

data "authentik_policy_binding" "policy-test2-30" {
  policy = authentik_policy_binding.policy-test2-30.policy
  target = authentik_policy_binding.policy-test2-30.target
  order  = authentik_policy_binding.policy-test2-30.order
}
`, appName, policyName)
}

func testAccDataSourcePolicyBindingUserConfig(appName string, userName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_user" "test1" {
  username = "%[2]s"
}

resource "authentik_user" "test2" {
  username = "%[2]s-2"
}

resource "authentik_policy_binding" "user-test1" {
  user   = authentik_user.test1.id
  target = authentik_application.test.uuid
  order  = 10
}

resource "authentik_policy_binding" "user-test2-20" {
  user   = authentik_user.test2.id
  target = authentik_application.test.uuid
  order  = 20
}

resource "authentik_policy_binding" "user-test2-30" {
  user   = authentik_user.test2.id
  target = authentik_application.test.uuid
  order  = 30
}

data "authentik_policy_binding" "user-test1" {
  user   = authentik_policy_binding.user-test1.user
  target = authentik_policy_binding.user-test1.target
  order  = authentik_policy_binding.user-test1.order
}

data "authentik_policy_binding" "user-test2-20" {
  user   = authentik_policy_binding.user-test2-20.user
  target = authentik_policy_binding.user-test2-20.target
  order  = authentik_policy_binding.user-test2-20.order
}

data "authentik_policy_binding" "user-test2-30" {
  user   = authentik_policy_binding.user-test2-30.user
  target = authentik_policy_binding.user-test2-30.target
  order  = authentik_policy_binding.user-test2-30.order
}
`, appName, userName)
}

const testAccDataSourcePolicyBindingNoConfigSimple = `
data "authentik_policy_binding" "missing-config" {
}
`

const testAccDataSourcePolicyBindingBrokenConfig1Simple = `
data "authentik_policy_binding" "broken-config-1" {
  id     = "12345678-1234-1234-1234-123456789012"
  target = "12345678-1234-1234-1234-123456789012"
}
`

const testAccDataSourcePolicyBindingBrokenConfig2Simple = `
data "authentik_policy_binding" "broken-config-2" {
  target = "12345678-1234-1234-1234-123456789012"
}
`

func testAccDataSourcePolicyBindingNotFound(appName string, groupName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_group" "test" {
  name = "%[2]s"
}

data "authentik_policy_binding" "group-test" {
  group  = authentik_group.test.id
  target = authentik_application.test.uuid
}
`, appName, groupName)
}

func testAccDataSourcePolicyBindingMultiple(appName string, groupName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_group" "test" {
  name = "%[2]s"
}

resource "authentik_policy_binding" "group-test-10" {
  group  = authentik_group.test.id
  target = authentik_application.test.uuid
  order  = 10
}

resource "authentik_policy_binding" "group-test-20" {
  group  = authentik_group.test.id
  target = authentik_application.test.uuid
  order  = 20
}

data "authentik_policy_binding" "duplicate-test" {
  group  = authentik_policy_binding.group-test-10.group
  target = authentik_policy_binding.group-test-20.target
}
`, appName, groupName)
}
