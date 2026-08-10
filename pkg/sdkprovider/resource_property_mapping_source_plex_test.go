package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePropertyMappingSourcePlex(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingSourcePlex(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_plex.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingSourcePlex(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_source_plex.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingSourcePlex(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_source_plex" "name" {
  name         = "%[1]s"
  expression   = "return True"
}
`, name)
}
