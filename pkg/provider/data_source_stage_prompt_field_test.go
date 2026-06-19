package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceStagePromptField(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStagePromptFieldByNameSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_stage_prompt_field.by-name", "name", "default-user-settings-field-email"),
				),
			},
			{
				Config: testAccDataSourceStagePromptFieldByIdSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_stage_prompt_field.by-id", "name", "default-user-settings-field-email"),
				),
			},
		},
	})
}

func TestAccDataSourceStagePromptField_Errors(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceStagePromptFieldMissingConfigSimple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourceStagePromptFieldBrokenConfigSimple,
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			},
			{
				Config:      testAccDataSourceStagePromptFieldNotFoundSimple,
				ExpectError: regexp.MustCompile(`No matching stage prompt fields found`),
			},
		},
	})
}

const testAccDataSourceStagePromptFieldByIdSimple = `
data "authentik_stage_prompt_field" "default_user_settings_field_email" {
  name = "default-user-settings-field-email"
}

data "authentik_stage_prompt_field" "by-id" {
  id = data.authentik_stage_prompt_field.default_user_settings_field_email.id
}
`

const testAccDataSourceStagePromptFieldByNameSimple = `
data "authentik_stage_prompt_field" "by-name" {
  name = "default-user-settings-field-email"
}
`

const testAccDataSourceStagePromptFieldMissingConfigSimple = `
data "authentik_stage_prompt_field" "missing-config" {
}
`

const testAccDataSourceStagePromptFieldBrokenConfigSimple = `
data "authentik_stage_prompt_field" "broken-config" {
  id   = "12345678-1234-1234-1234-123456789012"
  name = "probably-doesnt-exist-by-default"
}
`
const testAccDataSourceStagePromptFieldNotFoundSimple = `
data "authentik_stage_prompt_field" "not-found" {
  name = "probably-doesnt-exist-by-default"
}
`
