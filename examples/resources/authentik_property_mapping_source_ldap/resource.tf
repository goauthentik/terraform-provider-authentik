# Create a custom LDAP source property mapping

resource "authentik_property_mapping_source_ldap" "name" {
  name       = "custom-field"
  expression = "return ldap.get('sAMAccountName')"
}
