package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePropertyMappingProviderScope(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingProviderScope(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scope.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingProviderScope(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scope.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingProviderScope(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_provider_scope" "name" {
  name       = "%[1]s"
  scope_name = "%[1]s"
  expression = "return True"
}
`, name)
}
