//go:build enterprise

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceStageEndpoints(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageEndpoints(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_endpoints.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageEndpoints(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_endpoints.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageEndpoints(name string) string {
	return fmt.Sprintf(`
resource "authentik_endpoints_connector_agent" "name" {
  name = "%[1]s"
}

resource "authentik_stage_endpoints" "name" {
  name = "%[1]s"
  connector = authentik_endpoints_connector_agent.name.id
}
`, name)
}
