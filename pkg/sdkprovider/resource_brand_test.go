package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceBrand(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
	// TODO: branding_logo should be optional
	// TODO: branding_favicon should be optional
	return fmt.Sprintf(`
resource "authentik_brand" "name" {
  domain = "%[1]s"
  branding_logo = "test"
  branding_favicon = "test"
}
`, name)
}
