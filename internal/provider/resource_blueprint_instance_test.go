package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBlueprintInstance(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBlueprintInstanceSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_blueprint.instance", "name", rName),
				),
			},
			{
				Config: testAccResourceBlueprintInstanceSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_blueprint.instance", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceBlueprintInstanceSimple(name string) string {
	return fmt.Sprintf(`
resource "authentik_blueprint" "instance" {
  name = "%[1]s"
  path = "default/flow-default-authentication-flow.yaml"
}
`, name)
}
