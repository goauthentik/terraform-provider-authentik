# To get the ID of a SAML Property mapping

data "authentik_property_mapping_saml" "test" {
  managed = "goauthentik.io/providers/saml/upn"
}

# Then use `data.authentik_property_mapping_saml.admins.id`
