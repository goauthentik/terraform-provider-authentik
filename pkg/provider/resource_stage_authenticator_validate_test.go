package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func TestResourceStageAuthenticatorValidateSchemaToProvider(t *testing.T) {
	for _, tc := range []struct {
		name     string
		action   string
		expected api.NotConfiguredActionEnum
	}{
		{"skip", "skip", api.NOTCONFIGUREDACTIONENUM_SKIP},
		{"deny", "deny", api.NOTCONFIGUREDACTIONENUM_DENY},
		{"configure", "configure", api.NOTCONFIGUREDACTIONENUM_CONFIGURE},
		// `not_configured_action` is Required, so it must always end up in the request
		// body. Sourcing it via `helpers.GetP` (`d.GetOk`) yielded a nil pointer for the
		// zero value, and the model's `omitempty` then dropped the field from the UPDATE
		// PUT entirely, leaving authentik on its own default.
		{"zero value is still serialized", "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceStageAuthenticatorValidate().Schema, map[string]any{
				"name":                  "test",
				"not_configured_action": tc.action,
			})
			r := resourceStageAuthenticatorValidateSchemaToProvider(d)

			require.NotNil(t, r.NotConfiguredAction, "not_configured_action must never be nil")
			assert.Equal(t, tc.expected, *r.NotConfiguredAction)

			b, err := json.Marshal(r)
			assert.NoError(t, err)
			assert.Contains(t, string(b), `"not_configured_action"`, "not_configured_action must never be omitted from the request body")
		})
	}
}

func TestAccResourceStageAuthenticatorValidate(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
