package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceFlowStageBinding(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFlowStageBindingSimple(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName),
				),
			},
			{
				Config: testAccResourceFlowStageBindingSimple(rName+"test", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_flow.flow", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceFlowStageBindingSimple(name string, order int) string {
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
  order = %[2]d
}
`, name, order)
}
