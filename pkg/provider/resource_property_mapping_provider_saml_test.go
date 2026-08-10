package provider_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
)

func TestAccResourcePropertyMappingProviderSAML(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckPropertyMappingProviderSAMLDestroy,
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
			// Deliberately placed while expression is still the single-line "return True"
			// with no trailing newline, so this is an exact round trip. Against the heredoc
			// config below it would fail, and legitimately so: import starts from an empty
			// model, so ExpressionOrNull sees a null prior and state ends up holding the
			// API's stripped value rather than the config's newline-terminated one. That is
			// the documented consequence of replacing DiffSuppressFunc with semantic
			// equality, not a bug in this resource.
			{
				ResourceName:      "authentik_property_mapping_provider_saml.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The reason this resource carries the batch's extra coverage: a heredoc
			// expression ends in "\n" while authentik always returns it stripped, which is
			// exactly what helpers.ExpressionType's semantic equality exists to absorb.
			// No explicit plan check is needed - terraform-plugin-testing already runs a
			// plan after every apply and fails the step if it is non-empty, so this step
			// passing *is* the assertion that the trailing newline does not produce a
			// permanent diff.
			{
				Config: testAccResourcePropertyMappingProviderSAMLHeredoc(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_property_mapping_provider_saml.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccCheckPropertyMappingProviderSAMLDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_property_mapping_provider_saml" {
			continue
		}
		_, hr, err := c.PropertymappingsAPI.PropertymappingsProviderSamlRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("SAML provider property mapping %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
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

func testAccResourcePropertyMappingProviderSAMLHeredoc(name string) string {
	return fmt.Sprintf(`
resource "authentik_property_mapping_provider_saml" "name" {
  name      = "%[1]s"
  saml_name = "%[1]s"
  expression = <<-EOT
    return True
  EOT
}
`, name)
}
