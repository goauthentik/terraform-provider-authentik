package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderOAuth2DCR(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2DCR(rName, "hours=1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2_dcr.dcr", "access_token_validity", "hours=1"),
					resource.TestCheckResourceAttrSet("authentik_provider_oauth2_dcr.dcr", "oauth2_provider"),
				),
			},
			{
				Config: testAccResourceProviderOAuth2DCR(rName, "hours=2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2_dcr.dcr", "access_token_validity", "hours=2"),
				),
			},
		},
	})
}

func testAccResourceProviderOAuth2DCR(name string, accessTokenValidity string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name              = "authentik Self-signed Certificate"
  fetch_key         = false
  fetch_certificate = false
}

resource "authentik_provider_oauth2" "name" {
  name                = "%[1]s"
  client_id           = "%[1]s"
  signing_key         = data.authentik_certificate_key_pair.generated.id
  authorization_flow  = data.authentik_flow.default-authorization-flow.id
  invalidation_flow   = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_provider_oauth2_dcr" "dcr" {
  oauth2_provider        = authentik_provider_oauth2.name.id
  access_token_validity  = "%[2]s"
  allowed_grant_types    = ["authorization_code", "refresh_token"]
}
`, name, accessTokenValidity)
}
