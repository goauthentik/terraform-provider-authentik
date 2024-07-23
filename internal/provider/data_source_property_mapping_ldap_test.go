package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLDAPPropertyMapping(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLDAPPropertyMappingSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "name", "authentik default LDAP Mapping: Name"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "managed", "goauthentik.io/sources/ldap/default-name"),
				),
			},
			{
				Config: testAccDataSourceLDAPPropertyMappingList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourceLDAPPropertyMappingSimple = `
data "authentik_property_mapping_ldap" "test" {
  name    = "authentik default LDAP Mapping: Name"
  managed = "goauthentik.io/sources/ldap/default-name"
}
`

const testAccDataSourceLDAPPropertyMappingList = `
data "authentik_property_mapping_ldap" "test" {
  managed_list = [
    "goauthentik.io/sources/ldap/default-name",
    "goauthentik.io/sources/ldap/default-mail"
  ]
}
`
