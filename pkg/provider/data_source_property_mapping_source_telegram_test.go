package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePropertyMappingSourceTelegram(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropertyMappingSourceTelegram(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_source_telegram.test", "name", rName),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_source_telegram.test", "expression", "return True"),
				),
			},
		},
	})
}

func testAccDataSourcePropertyMappingSourceTelegram(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_source_telegram" "test" {
  name       = "%[1]s"
  expression = "return True"
}

data "authentik_property_mapping_source_telegram" "test" {
  name = "%[1]s"
}
`, name)
}
