# To get the ID of a SCIM Property mapping

data "authentik_property_mapping_scim" "test" {
  managed = "goauthentik.io/providers/scim/user"
}

# Then use `data.authentik_property_mapping_scim.test.id`

# Or, to get the IDs of multiple mappings

data "authentik_property_mapping_scim" "test" {
  managed_list = [
    "goauthentik.io/providers/scim/user",
    "goauthentik.io/providers/scim/group"
  ]
}

# Then use data.authentik_property_mapping_scim.test.ids
