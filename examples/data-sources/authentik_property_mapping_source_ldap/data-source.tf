# To get the ID of a LDAP Source Property mapping

data "authentik_property_mapping_source_ldap" "test" {
  managed = "goauthentik.io/sources/ldap/default-name"
}

# Then use `data.authentik_property_mapping_source_ldap.test.id`

# Or, to get the IDs of multiple mappings

data "authentik_property_mapping_source_ldap" "test" {
  managed_list = [
    "goauthentik.io/sources/ldap/default-name",
    "goauthentik.io/sources/ldap/default-mail"
  ]
}

# Then use data.authentik_property_mapping_source_ldap.test.ids
