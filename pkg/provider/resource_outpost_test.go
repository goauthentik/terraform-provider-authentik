package provider

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "protocol_providers.#", "2"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "type", "proxy"),
				),
			},
			{
				Config: testAccResourceOutpostSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "protocol_providers.#", "2"),
					resource.TestCheckResourceAttr("authentik_outpost.outpost", "type", "proxy"),
				),
			},
			{
				Config:      testAccResourceOutpostInvalidConfig(rName + "test"),
				ExpectError: regexp.MustCompile("invalid character"),
			},
		},
	})
}

func testAccResourceOutpostSimple(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_proxy" "proxy" {
  name               = "%[1]s-1"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
}

resource "authentik_provider_proxy" "proxy-2" {
  name               = "%[1]s-2"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
}

resource "authentik_outpost" "outpost" {
  name = "%[1]s"
  protocol_providers = [
    authentik_provider_proxy.proxy.id,
    authentik_provider_proxy.proxy-2.id,
  ]
  config = jsonencode(
    {
      authentik_host                 = "http://localhost:9000/"
      authentik_host_browser         = ""
    }
  )
}
`, name)
}

func testAccResourceOutpostInvalidConfig(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_proxy" "proxy" {
  name               = "proxy"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
}

resource "authentik_outpost" "outpost" {
  name = "%[1]s"
  protocol_providers = [
    authentik_provider_proxy.proxy.id
  ]
  config = "a"
}
`, name)
}
