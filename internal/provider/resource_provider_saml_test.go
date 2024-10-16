package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderSAML(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderSAML(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_saml.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderSAML(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_saml.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderSAML(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_property_mapping_provider_saml" "test" {
  managed_list = [
    "goauthentik.io/providers/saml/upn",
    "goauthentik.io/providers/saml/name"
  ]
}
resource "authentik_property_mapping_provider_saml" "name" {
  name       = "%[1]s"
  saml_name  = "%[1]s"
  expression = "return True"
}

resource "authentik_provider_saml" "name" {
  name               = "%[1]s"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow  = data.authentik_flow.default-provider-invalidation-flow.id
  acs_url            = "http://localhost"
  property_mappings  = concat(data.authentik_property_mapping_provider_saml.test.ids, [
    authentik_property_mapping_provider_saml.name.id
  ])
}
resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_saml.name.id
}
`, name, appName)
}
