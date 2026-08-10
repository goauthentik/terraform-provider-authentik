package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceProviderOAuth2(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckProviderOAuth2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.0.redirect_uri_type", "authorization"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.1.redirect_uri_type", "logout"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderOAuth2(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.0.redirect_uri_type", "authorization"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.1.redirect_uri_type", "logout"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
			{
				ResourceName:      "authentik_provider_oauth2.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckProviderOAuth2Destroy(s *terraform.State) error {
	c := testAccAPIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_provider_oauth2" {
			continue
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}
		_, hr, err := c.ProvidersAPI.ProvidersOauth2Retrieve(context.Background(), int32(id)).Execute()
		if err == nil {
			return fmt.Errorf("provider_oauth2 %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func TestAccResourceProviderOAuth2_WithSecret(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2WithSecret(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderOAuth2WithSecret(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderOAuth2(name string, appName string) string {
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
  client_id = "%[1]s"
  # client_secret = "test"
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  allowed_redirect_uris = [
    {
      matching_mode     = "strict",
      url               = "http://localhost/callback",
      redirect_uri_type = "authorization",
    },
    {
      matching_mode     = "strict",
      url               = "http://localhost/",
      redirect_uri_type = "logout",
    }
  ]
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_oauth2.name.id
}
`, name, appName)
}

func testAccResourceProviderOAuth2WithSecret(name string, appName string) string {
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
  client_id = "%[1]s"
  client_secret = "test"
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_oauth2.name.id
}
`, name, appName)
}
