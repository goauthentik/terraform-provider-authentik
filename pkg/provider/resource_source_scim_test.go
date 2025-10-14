package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSourceSCIM(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceSCIM(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_scim.name", "name", rName),
					resource.TestCheckResourceAttrSet("authentik_source_scim.name", "scim_url"),
					resource.TestCheckResourceAttrSet("authentik_source_scim.name", "token"),
				),
			},
			{
				Config: testAccResourceSourceSCIM(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_scim.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSourceSCIM(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_scim" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
}
`, name, appName)
}
