# To get the ID of a SAML Property mapping

data "authentik_property_mapping_saml" "test" {
  managed = "goauthentik.io/providers/saml/upn"
}

# Then use `data.authentik_property_mapping_saml.test.id`

# Or, to get the IDs of multiple mappings

data "authentik_property_mapping_saml" "test" {
  managed_list = [
    "goauthentik.io/providers/saml/upn",
    "goauthentik.io/providers/saml/name"
  ]
}

# Then use data.authentik_property_mapping_saml.test.ids
