package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourcePropertyMappingSourceSCIM(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingSourceSCIM(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_scim.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingSourceSCIM(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_scim.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingSourceSCIM(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_source_scim" "name" {
  name         = "%[1]s"
  expression   = "return True"
}
`, name)
}
