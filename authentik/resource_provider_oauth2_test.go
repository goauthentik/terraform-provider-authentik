package authentik

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderOAuth2(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", "grafana"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", "grafana"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", "test app"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", "test-app"),
				),
			},
		},
	})
}

const testAccResourceProviderOAuth2 = `
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_oauth2" "name" {
  name      = "grafana"
  client_id = "grafana"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
}

resource "authentik_application" "name" {
  name              = "test app"
  slug              = "test-app"
  protocol_provider = authentik_provider_oauth2.name.id
}
`
