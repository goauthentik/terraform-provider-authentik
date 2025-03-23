package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSourceLDAP(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourceLDAP(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_ldap.name", "name", rName),
				),
			},
			{
				Config: testAccResourceSourceLDAP(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_ldap.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceSourceLDAP(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_property_mapping_source_ldap" "test" {
  managed_list = [
    "goauthentik.io/sources/ldap/default-name",
    "goauthentik.io/sources/ldap/default-mail"
  ]
}

resource "authentik_source_ldap" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"

  server_uri = "ldaps://1.2.3.4"
  bind_cn = "foo"
  bind_password = "bar"
  base_dn = "dn=foo"
  sync_users_password = false
  property_mappings = data.authentik_property_mapping_source_ldap.test.ids
  property_mappings_group = data.authentik_property_mapping_source_ldap.test.ids
}
`, name, appName)
}
