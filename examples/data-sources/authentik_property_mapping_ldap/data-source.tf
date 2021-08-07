# To get the ID of a LDAP Property mapping

data "authentik_property_mapping_ldap" "test" {
  managed = "goauthentik.io/sources/ldap/default-name"
}

# Then use `data.authentik_property_mapping_ldap.test.id`
