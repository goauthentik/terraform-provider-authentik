package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePropertyMappingSourceLDAP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePropertyMappingSourceLDAPSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_source_ldap.test", "name", "authentik default LDAP Mapping: Name"),
					resource.TestCheckResourceAttr("data.authentik_property_mapping_source_ldap.test", "managed", "goauthentik.io/sources/ldap/default-name"),
				),
			},
			{
				Config: testAccDataSourcePropertyMappingSourceLDAPList,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_property_mapping_source_ldap.test", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccDataSourcePropertyMappingSourceLDAPSimple = `
data "authentik_property_mapping_source_ldap" "test" {
  name    = "authentik default LDAP Mapping: Name"
  managed = "goauthentik.io/sources/ldap/default-name"
}
`

const testAccDataSourcePropertyMappingSourceLDAPList = `
data "authentik_property_mapping_source_ldap" "test" {
  managed_list = [
    "goauthentik.io/sources/ldap/default-name",
    "goauthentik.io/sources/ldap/default-mail"
  ]
}
`
