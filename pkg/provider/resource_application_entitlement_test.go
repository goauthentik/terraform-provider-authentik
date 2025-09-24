package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceApplicationEntitlement(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationEntitlementSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName),
					resource.TestCheckResourceAttr("authentik_application_entitlement.ent", "name", rName),
				),
			},
			{
				Config: testAccResourceApplicationEntitlementSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application_entitlement.ent", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceApplicationEntitlementSimple(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authentication-flow" {
  slug = "default-authentication-flow"
}

data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
}

resource "authentik_application" "name" {
  name              = "%[1]s"
  slug              = "%[1]s"
  meta_icon = "http://localhost/%[1]s"
}

resource "authentik_application_entitlement" "ent" {
  name = "%[1]s"
  application = authentik_application.name.uuid
}
`, name)
}
