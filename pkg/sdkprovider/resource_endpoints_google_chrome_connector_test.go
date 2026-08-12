package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceEndpointsGoogleChromeConnector(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
