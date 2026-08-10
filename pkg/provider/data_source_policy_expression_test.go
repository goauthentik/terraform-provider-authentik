package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePolicyExpression(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyExpressionByNameSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_expression.by-name", "name", "default-user-settings-authorization"),
				),
			},
			{
				Config: testAccDataSourcePolicyExpressionByIdSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_policy_expression.by-id", "name", "default-user-settings-authorization"),
				),
			},
		},
	})
}

func TestAccDataSourcePolicyExpression_Errors(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourcePolicyExpressionMissingConfigSimple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourcePolicyExpressionBrokenConfigSimple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourcePolicyExpressionNotFoundSimple,
				ExpectError: regexp.MustCompile(`No matching expression policies found`),
			},
		},
	})
}

const testAccDataSourcePolicyExpressionByIdSimple = `
data "authentik_policy_expression" "default_user_settings_authorization" {
  name = "default-user-settings-authorization"
}

data "authentik_policy_expression" "by-id" {
  id = data.authentik_policy_expression.default_user_settings_authorization.id
}
`

const testAccDataSourcePolicyExpressionByNameSimple = `
data "authentik_policy_expression" "by-name" {
  name = "default-user-settings-authorization"
}
`

const testAccDataSourcePolicyExpressionMissingConfigSimple = `
data "authentik_policy_expression" "missing-config" {
}
`

const testAccDataSourcePolicyExpressionBrokenConfigSimple = `
data "authentik_policy_expression" "broken-config" {
  id   = "12345678-1234-1234-1234-123456789012"
  name = "probably-doesnt-exist-by-default"
}
`
const testAccDataSourcePolicyExpressionNotFoundSimple = `
data "authentik_policy_expression" "not-found" {
  name = "probably-doesnt-exist-by-default"
}
`
