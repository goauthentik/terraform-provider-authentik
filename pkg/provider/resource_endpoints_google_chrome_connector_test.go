//go:build enterprise

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEndpointsGoogleChromeConnector(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEndpointsGoogleChromeConnector(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_google_chrome_connector.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_endpoints_google_chrome_connector.name", "enabled", "true"),
				),
			},
			{
				Config: testAccResourceEndpointsGoogleChromeConnector(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_google_chrome_connector.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceEndpointsGoogleChromeConnector(name string) string {
	return fmt.Sprintf(`
resource "authentik_endpoints_google_chrome_connector" "name" {
  name        = "%[1]s"
  credentials = jsonencode({
    type         = "service_account"
    project_id   = "example-project"
    client_email = "example@example-project.iam.gserviceaccount.com"
  })
}
`, name)
}
