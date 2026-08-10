package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePropertyMappingProviderScope(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingProviderScope(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scope.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingProviderScope(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_scope.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingProviderScope(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_provider_scope" "name" {
  name       = "%[1]s"
  scope_name = "%[1]s"
  expression = "return True"
}
`, name)
}
