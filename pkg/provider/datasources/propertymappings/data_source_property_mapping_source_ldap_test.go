package propertymappings_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestAccDataSourcePropertyMappingSourceLDAP(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
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
