package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePropertyMappingNotification(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingNotification(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_notification.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingNotification(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_notification.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingNotification(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_notification" "name" {
  name       = "%[1]s"
  expression = "return True"
}
`, name)
}
