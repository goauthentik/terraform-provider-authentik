package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderRadius(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderRadius(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_radius.name", "name", rName),
				),
			},
			{
				Config: testAccResourceProviderRadius(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_radius.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderRadius(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authentication-flow" {
  slug = "default-authentication-flow"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_radius" "name" {
  name      = "%[1]s"
  authorization_flow = data.authentik_flow.default-authentication-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  shared_secret = "foo"
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_radius.name.id
}
`, name, appName)
}
