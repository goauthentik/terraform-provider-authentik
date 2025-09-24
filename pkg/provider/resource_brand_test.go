package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBrand(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBrand(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_brand.name", "domain", rName),
				),
			},
			{
				Config: testAccResourceBrand(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_brand.name", "domain", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceBrand(name string) string {
	return fmt.Sprintf(`
# TODO: branding_logo should be optional
# TODO: branding_favicon should be optional
resource "authentik_brand" "name" {
  domain = "%[1]s"
  branding_logo = "test"
  branding_favicon = "test"
}
`, name)
}
