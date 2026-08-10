package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageAuthenticatorSms(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorSms(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_sms.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorSms(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_sms.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorSms(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_sms" "name" {
  name              = "%[1]s"
  from_number = "1234"
  account_sid = "bar"
  auth = "baz"
}
`, name)
}
