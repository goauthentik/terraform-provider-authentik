package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScopeMapping(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScopeMapping(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_scope_mapping.name", "name", rName),
				),
			},
			{
				Config: testAccResourceScopeMapping(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_scope_mapping.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceScopeMapping(name string) string {
	return fmt.Sprintf(`
resource "authentik_scope_mapping" "name" {
  name       = "%[1]s"
  scope_name = "%[1]s"
  expression = "return True"
}
`, name)
}
