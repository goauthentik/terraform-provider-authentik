package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutpost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutpostConfig("test-outpost-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_outpost.test", "name", "test-outpost-ds"),
				),
			},
		},
	})
}

func testAccDataSourceOutpostConfig(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_proxy" "test" {
  name               = "%s-proxy"
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow  = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_outpost" "test" {
  name = "%s"
  protocol_providers = [authentik_provider_proxy.test.id]
}

data "authentik_outpost" "test" {
  name = authentik_outpost.test.name
}
`, name, name)
}
