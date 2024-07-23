# Create a custom LDAP property mapping

resource "authentik_property_mapping_ldap" "name" {
  name       = "custom-field"
  expression = "return ldap.get('sAMAccountName')"
}
