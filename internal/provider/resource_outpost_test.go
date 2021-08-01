package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOutpost(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOutpostSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "name", rName),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "protocol_providers.#", "1"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "type", "proxy"),
				),
			},
		},
	})
}

func testAccResourceOutpostSimple(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_proxy" "proxy" {
  name               = "proxy"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
}

resource "authentik_outpost" "outpost" {
  name = "%s"
  protocol_providers = [
    authentik_provider_proxy.proxy.id
  ]
}
`, name)
}
