package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageAuthenticatorValidate(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorValidate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorValidate(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorValidate(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_totp" "name" {
  name              = "%[1]s-setup"
}

resource "authentik_stage_authenticator_validate" "name" {
  name              = "%[1]s"
  device_classes = ["static"]
  not_configured_action = "skip"
  configuration_stages = [
    authentik_stage_authenticator_totp.name.id,
  ]
}
`, name)
}
