package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTenant(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTenant(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_tenant.name", "domain", rName),
				),
			},
			{
				Config: testAccResourceTenant(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_tenant.name", "domain", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceTenant(name string) string {
	return fmt.Sprintf(`
# TODO: branding_logo should be optional
# TODO: branding_favicon should be optional
resource "authentik_tenant" "name" {
  domain = "%s"
  branding_logo = "test"
  branding_favicon = "test"
}
`, name)
}
