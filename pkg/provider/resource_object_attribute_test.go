package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceObjectAttribute(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceObjectAttribute(rName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_object_attribute.name", "key", rName),
					resource.TestCheckResourceAttr("authentik_object_attribute.name", "is_required", "false"),
				),
			},
			{
				Config: testAccResourceObjectAttribute(rName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_object_attribute.name", "key", rName),
					resource.TestCheckResourceAttr("authentik_object_attribute.name", "is_required", "true"),
				),
			},
		},
	})
}

func testAccResourceObjectAttribute(name string, isRequired string) string {
	return fmt.Sprintf(`
resource "authentik_object_attribute" "name" {
  object_type = "authentik_core.user"
  key         = "%[1]s"
  label       = "%[1]s"
  type        = "text"
  is_required = %[2]s
}
`, name, isRequired)
}
