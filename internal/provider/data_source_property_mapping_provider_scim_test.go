package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePropertyMappingProviderSCIM(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropertyMappingProviderSCIMSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scim.test", "name", "authentik default SCIM Mapping: User"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scim.test", "managed", "goauthentik.io/providers/scim/user"),
				),
			},
			{
				Config: testAccDataSourcePropertyMappingProviderSCIMList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_provider_scim.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourcePropertyMappingProviderSCIMSimple = `
data "authentik_property_mapping_provider_scim" "test" {
  name    = "authentik default SCIM Mapping: User"
  managed = "goauthentik.io/providers/scim/user"
}
`

const testAccDataSourcePropertyMappingProviderSCIMList = `
data "authentik_property_mapping_provider_scim" "test" {
  managed_list = [
    "goauthentik.io/providers/scim/user",
    "goauthentik.io/providers/scim/group"
  ]
}
`
