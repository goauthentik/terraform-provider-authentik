package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePropertyMappingProviderSAML(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePropertyMappingProviderSAML(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_saml.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePropertyMappingProviderSAML(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_saml.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePropertyMappingProviderSAML(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_provider_saml" "name" {
  name       = "%[1]s"
  saml_name  = "%[1]s"
  expression = "return True"
}
`, name)
}
