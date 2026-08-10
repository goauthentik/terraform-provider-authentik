package sdkprovider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
)

func TestAccResourceStageAuthenticatorValidate(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "skip"),
				Check:  testAccResourceStageAuthenticatorValidateCheck(rName, "skip"),
			},
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "deny"),
				Check:  testAccResourceStageAuthenticatorValidateCheck(rName, "deny"),
			},
			{
				Config: testAccResourceStageAuthenticatorValidateAction(rName, "configure"),
				Check:  testAccResourceStageAuthenticatorValidateCheck(rName, "configure"),
			},
		},
	})
}

// testAccResourceStageAuthenticatorValidateCheck asserts every attribute, not just the one being
// flipped, so that a field silently dropped from the UPDATE body shows up. All values below
// deliberately differ from the schema/server defaults, otherwise a dropped field is invisible.
func testAccResourceStageAuthenticatorValidateCheck(name string, action string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "name", name),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "not_configured_action", action),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "device_classes.#", "1"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "device_classes.0", "static"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "configuration_stages.#", "1"),
		resource.TestCheckResourceAttrPair(
			"authentik_stage_authenticator_validate.name", "configuration_stages.0",
			"authentik_stage_authenticator_totp.name", "id",
		),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "last_auth_threshold", "minutes=5"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "webauthn_user_verification", "required"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "webauthn_hints.#", "1"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "webauthn_hints.0", "security-key"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "email_otp_throttling_factor", "2"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "sms_otp_throttling_factor", "2"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "totp_otp_throttling_factor", "2"),
		resource.TestCheckResourceAttr("authentik_stage_authenticator_validate.name", "static_otp_throttling_factor", "2"),
	)
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
  last_auth_threshold = "minutes=5"
  webauthn_user_verification = "required"
  webauthn_hints = ["security-key"]
  email_otp_throttling_factor = 2
  sms_otp_throttling_factor = 2
  totp_otp_throttling_factor = 2
  static_otp_throttling_factor = 2
}
`, name, action)
}
