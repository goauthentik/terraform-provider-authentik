package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	api "goauthentik.io/api/v3"
)

func TestAccResourceApplication(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "protocol_provider", "0"),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_launch_url", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_description", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_publisher", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "policy_engine_mode", string(api.POLICYENGINEMODE_ANY)),
				),
			},
			{
				Config: testAccResourceApplicationSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "protocol_provider", "0"),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_launch_url", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_description", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_publisher", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "policy_engine_mode", string(api.POLICYENGINEMODE_ANY)),
				),
			},
			{
				Config:      testAccResourceApplicationSimple(rName + "test+"),
				ExpectError: regexp.MustCompile("consisting of letters, numbers, underscores or hyphens"),
			},
		},
	})
}

func testAccResourceApplicationSimple(name string) string {
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

resource "authentik_provider_ldap" "name" {
  name      = "%[1]s"
  base_dn = "dc=%[1]s,dc=goauthentik,dc=io"
  bind_flow = data.authentik_flow.default-authentication-flow.id
  unbind_flow = data.authentik_flow.default-provider-invalidation-flow.id
  tls_server_name = "foo"
  certificate = data.authentik_certificate_key_pair.generated.id
}

resource "authentik_application" "name" {
  name              = "%[1]s"
  slug              = "%[1]s"
  meta_icon = "http://localhost/%[1]s"
  backchannel_providers = [authentik_provider_ldap.name.id]
}
`, name)
}
