package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSourcePlex(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourcePlex(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_plex.name", "name", rName),
				),
			},
			{
				Config: testAccResourceSourcePlex(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_plex.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSourcePlex(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_plex" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id
  client_id = "foo-bar-baz"
  plex_token = "foo"
}
`, name, appName)
}
