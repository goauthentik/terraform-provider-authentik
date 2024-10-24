//go:build enterprise
// +build enterprise

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageAuthenticatorEndpointGDTC(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorEndpointGDTC(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_endpoint_gdtc.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorEndpointGDTC(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_endpoint_gdtc.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorEndpointGDTC(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_endpoint_gdtc" "name" {
  name              = "%[1]s"
  credentials = jsonencode({
    "type"="service_account",
    "project_id"="foo",
    "private_key_id"="foo",
    "private_key"="bar\n",
    "client_email"="bar",
    "client_id"="qewrqer",
    "auth_uri"="https://accounts.google.com/o/oauth2/auth",
    "token_uri"="https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url"="https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url"="foo",
    "universe_domain"="googleapis.com"
  })
}
`, name)
}
