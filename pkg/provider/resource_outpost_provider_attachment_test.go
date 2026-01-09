package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceOutpostProviderAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOutpostProviderAttachmentConfig("test-outpost"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("authentik_outpost_provider_attachment.test", "outpost"),
					resource.TestCheckResourceAttrSet("authentik_outpost_provider_attachment.test", "protocol_provider"),
				),
			},
		},
	})
}

func testAccResourceOutpostProviderAttachmentConfig(name string) string {
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

data "authentik_outpost" "emb" {
  name = "authentik Embedded Outpost"
}

resource "authentik_outpost_provider_attachment" "test" {
  outpost  = data.authentik_outpost.emb.id
  protocol_provider = authentik_provider_proxy.test.id
}
`, name)
}
