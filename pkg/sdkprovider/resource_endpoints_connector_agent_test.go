package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceEndpointsAgent(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEndpointsAgent(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_connector_agent.name", "name", rName),
				),
			},
			{
				Config: testAccResourceEndpointsAgent(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_connector_agent.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceEndpointsAgent(name string) string {
	return fmt.Sprintf(`
resource "authentik_endpoints_connector_agent" "name" {
  name = "%[1]s"
}
`, name)
}
