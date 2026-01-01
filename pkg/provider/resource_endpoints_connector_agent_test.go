package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEndpointsAgent(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
