package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSAMLPropertyMapping(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSAMLPropertyMappingSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_saml.test", "name", "authentik default SAML Mapping: UPN"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_saml.test", "managed", "goauthentik.io/providers/saml/upn"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_saml.test", "saml_name", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"),
				),
			},
		},
	})
}

const testAccDataSourceSAMLPropertyMappingSimple = `
data "authentik_property_mapping_saml" "test" {
  name    = "authentik default SAML Mapping: UPN"
  managed = "goauthentik.io/providers/saml/upn"
}
`
