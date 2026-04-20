package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceStageAuthenticatorValidate(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "skip"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "not_configured_action", "skip"),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "deny"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "not_configured_action", "deny"),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "configure"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "not_configured_action", "configure"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorValidateAction(name string, action string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_totp" "name" {
  name              = "%[1]s-setup"
}

resource "authentik_stage_authenticator_validate" "name" {
  name              = "%[1]s"
  device_classes = ["static"]
  not_configured_action = "%[2]s"
  configuration_stages = [
    authentik_stage_authenticator_totp.name.id,
  ]
}
`, name, action)
}
