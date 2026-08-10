package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderWSFederation(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderWSFederation(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_ws_federation.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderWSFederation(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_ws_federation.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderWSFederation(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_ws_federation" "name" {
  name               = "%[1]s"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow  = data.authentik_flow.default-provider-invalidation-flow.id
  reply_url          = "http://localhost"
  wtrealm            = "http://localhost"
}
resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_ws_federation.name.id
}
`, name, appName)
}
