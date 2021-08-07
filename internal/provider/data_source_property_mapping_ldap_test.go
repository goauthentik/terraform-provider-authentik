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
					resource.TestCheckResourceAttr("data.authentik_property_mapping_ldap.test", "object_field", "name"),
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
