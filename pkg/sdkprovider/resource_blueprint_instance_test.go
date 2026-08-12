package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceBlueprintInstance(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
