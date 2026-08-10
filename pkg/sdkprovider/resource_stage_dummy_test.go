package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageDummy(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageDummy(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_dummy.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageDummy(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_dummy.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageDummy(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_dummy" "name" {
  name              = "%[1]s"
}
`, name)
}
