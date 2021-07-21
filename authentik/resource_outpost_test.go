package authentik

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOutpost(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOutpostSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "name", "acc-test-app"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "protocol_providers.#", "1"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "type", "proxy"),
				),
			},
		},
	})
}

const testAccResourceOutpostSimple = `
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
  name = "acc-test-app"
  protocol_providers = [
    authentik_provider_proxy.proxy.id
  ]
}
`
