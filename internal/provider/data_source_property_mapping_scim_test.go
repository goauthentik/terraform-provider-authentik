package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSCIMPropertyMapping(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSCIMPropertyMappingSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "name", "authentik default SCIM Mapping: User"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "managed", "goauthentik.io/providers/scim/user"),
				),
			},
			{
				Config: testAccDataSourceSCIMPropertyMappingList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourceSCIMPropertyMappingSimple = `
data "authentik_property_mapping_scim" "test" {
  name    = "authentik default SCIM Mapping: User"
  managed = "goauthentik.io/providers/scim/user"
}
`

const testAccDataSourceSCIMPropertyMappingList = `
data "authentik_property_mapping_ldap" "test" {
  managed_list = [
    "goauthentik.io/providers/scim/user",
    "goauthentik.io/providers/scim/group
"
  ]
}
`
