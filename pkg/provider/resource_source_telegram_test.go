package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSourceTelegram(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceTelegram(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_telegram.name", "name", rName),
				),
			},
			{
				Config: testAccResourceSourceTelegram(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_telegram.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSourceTelegram(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_telegram" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id
  pre_authentication_flow = data.authentik_flow.default-authorization-flow.id
  bot_username = "foo"
  bot_token = "foo"
}
`, name, appName)
}
