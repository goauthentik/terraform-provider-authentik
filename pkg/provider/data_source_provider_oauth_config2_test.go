package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOAuth2ProviderConfig(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	clientId := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOAuth2ProviderConfigSimple(rName, clientId, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "issuer_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "authorize_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "token_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "user_info_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "provider_info_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "logout_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.name", "jwks_url"),
					resource.TestCheckResourceAttr("data.authentik_provider_oauth2_config.name", "client_id", clientId),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "issuer_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "authorize_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "token_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "user_info_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "provider_info_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "logout_url"),
					resource.TestCheckResourceAttrSet("data.authentik_provider_oauth2_config.id", "jwks_url"),
					resource.TestCheckResourceAttr("data.authentik_provider_oauth2_config.id", "client_id", clientId),
				),
			},
		},
	})
}

func testAccDataSourceOAuth2ProviderConfigSimple(name string, clientId string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}
data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
  fetch_key = false
  fetch_certificate = false
}
resource "authentik_provider_oauth2" "name" {
  name      = "%[1]s"
  client_id = "%[2]s"
  # client_secret = "test"
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_application" "name" {
  name              = "%[3]s"
  slug              = "%[3]s"
  protocol_provider = authentik_provider_oauth2.name.id
}

data "authentik_provider_oauth2_config" "name" {
  name = "%[1]s"
  depends_on = [
	authentik_application.name
  ]
}

data "authentik_provider_oauth2_config" "id" {
  provider_id = authentik_provider_oauth2.name.id
  depends_on = [
	authentik_application.name
  ]
}
`, name, clientId, appName)
}
