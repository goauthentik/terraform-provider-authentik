# Create a custom LDAP property mapping

resource "authentik_property_mapping_ldap" "name" {
  name         = "custom-field"
  object_field = "username"
  expression   = "return ldap.get('sAMAccountName')"
}
