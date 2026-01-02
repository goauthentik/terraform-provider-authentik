package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEndpointEnrollmenToken(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	expires := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEndpointEnrollmenToken(rName, expires),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_connector_agent_enrollment_token.token", "name", rName),
					resource.TestCheckResourceAttrSet("authentik_endpoints_connector_agent_enrollment_token.token", "key"),
				),
			},
		},
	})
}

func testAccResourceEndpointEnrollmenToken(name string, time string) string {
	return fmt.Sprintf(`
resource "authentik_endpoints_connector_agent" "agent" {
  name = "%[1]s"
}

resource "authentik_endpoints_connector_agent_enrollment_token" "token" {
	connector = authentik_endpoints_connector_agent.agent.id
	expires = "%[2]s"
	name = "%[1]s"
	retrieve_key = true
}
`, name, time)
}
