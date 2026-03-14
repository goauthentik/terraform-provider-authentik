package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceStageAuthenticatorWebAuthn(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorWebAuthn(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_webauthn.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorWebAuthn(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_webauthn.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorWebAuthn(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_webauthn" "name" {
  name              = "%[1]s"
}
`, name)
}
