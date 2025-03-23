package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePropertyMappingProviderSAML(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropertyMappingProviderSAMLSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_saml.test", "name", "authentik default SAML Mapping: UPN"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_saml.test", "managed", "goauthentik.io/providers/saml/upn"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_saml.test", "saml_name", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"),
				),
			},
			{
				Config: testAccDataSourcePropertyMappingProviderSAMLList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_saml.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourcePropertyMappingProviderSAMLSimple = `
data "authentik_property_mapping_provider_saml" "test" {
  name    = "authentik default SAML Mapping: UPN"
  managed = "goauthentik.io/providers/saml/upn"
}
`

const testAccDataSourcePropertyMappingProviderSAMLList = `
data "authentik_property_mapping_provider_saml" "test" {
  managed_list = [
    "goauthentik.io/providers/saml/upn",
    "goauthentik.io/providers/saml/name"
  ]
}
`
