package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStagePrompt(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStagePrompt(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_prompt.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStagePrompt(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_prompt.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStagePrompt(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_prompt_field" "field" {
  name = "%[1]s"
  field_key = "%[1]s-test-field"
  label = "a label"
  type = "text"
  placeholder            = <<-EOT
    try:
      return user.username
    except:
      return ''
EOT
}
resource "authentik_stage_prompt" "name" {
  name              = "%[1]s"
  fields = [
    resource.authentik_stage_prompt_field.field.id,
  ]
}
`, name)
}
