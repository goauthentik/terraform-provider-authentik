package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSourceOAuth(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceOAuth(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_oauth.name", "name", rName),
				),
			},
			{
				Config: testAccResourceSourceOAuth(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_oauth.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSourceOAuth(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_oauth" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id

  provider_type = "discord"
  consumer_key = "foo"
  consumer_secret = "bar"
}

resource "authentik_source_oauth" "name2" {
  name      = "%[1]s-jwks"
  slug      = "%[1]s-jwks"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id

  provider_type = "openidconnect"
  consumer_key = "foo"
  consumer_secret = "bar"

  access_token_url = "http://foo"
  authorization_url = "http://foo"
  profile_url = "http://foo"
  oidc_jwks = <<EOT
  {
	"keys": [
        {
            "use": "sig",
            "kty": "RSA",
            "kid": "zOsJ_KGpGiHYhVXZkyYNjRRi20RzHB_RGNqS-NNpg0U",
            "alg": "RS256",
            "n": "rl2d9iyUrGxQPerG869z3Jvs9x7v0RUTO8zb9hlx_nygTGDmaCLBEkmj5w89JCHDv8Imbp9IjklZtrEc6mIYfLghy6vzv95ouaRGzhtrnZxwZcxHxW5R_sEd-HVq_infgndBHhllcLdAYy-Y7MYcV4bS0VDGPh3KT29q6uGO1VPJ1y0eEosOVzeUi8CQE2xvmYQ_TonCP_7KOZQ0Q0CzQURDo-yhgl0ijTMvLQdc2wJwGdxW4XzsOPyzmuKGKSow65U2lkrg3JAyXrNlS7HqCgpVqo8A1zYbcxoP2fGu3zxhMw0ImrIpnKxd-U9ZY1GiOLokCXBFkH6YXa3MSW6nVQ",
            "e": "AQAB"
        }
    ]
  }
EOT
}
`, name, appName)
}
