package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePropertyMappingProviderSCIM(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingProviderSCIM(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scim.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingProviderSCIM(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scim.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingProviderSCIM(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_provider_scim" "name" {
  name       = "%[1]s"
  expression = "return True"
}
`, name)
}
