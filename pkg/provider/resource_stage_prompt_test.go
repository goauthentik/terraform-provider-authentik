package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStagePrompt(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
