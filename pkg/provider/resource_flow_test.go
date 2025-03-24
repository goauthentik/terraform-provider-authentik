package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFlow(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFlowSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName),
				),
			},
			{
				Config: testAccResourceFlowSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceFlowSimple(name string) string {
	return fmt.Sprintf(`
resource "authentik_flow" "flow" {
  name = "%[1]s"
  title ="%[1]s"
  slug  ="%[1]s"
  designation = "authorization"
  authentication = "none"
  background = "https://goauthentik.io"
  layout = "stacked"
}
`, name)
}
