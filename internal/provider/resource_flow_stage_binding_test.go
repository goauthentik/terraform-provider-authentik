package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFlowStageBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFlowStageBindingSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName),
				),
			},
			{
				Config: testAccResourceFlowStageBindingSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceFlowStageBindingSimple(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_dummy" "name" {
  name              = "%[1]s"
}

resource "authentik_flow" "flow" {
  name = "%[1]s"
  title ="%[1]s"
  slug  ="%[1]s"
  designation = "authorization"
}

resource "authentik_flow_stage_binding" "dummy-flow" {
  target = authentik_flow.flow.uuid
  stage = authentik_stage_dummy.name.id
  order = 0
}
`, name)
}
