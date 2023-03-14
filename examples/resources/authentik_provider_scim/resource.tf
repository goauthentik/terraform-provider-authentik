# Create a SCIM Provider

data "authentik_property_mapping_scim" "user" {
  managed = "goauthentik.io/providers/scim/user"
}

data "authentik_property_mapping_scim" "group" {
  managed = "goauthentik.io/providers/scim/group"
}

resource "authentik_provider_scim" "name" {
  name                    = "test-app"
  url                     = "http://localhost"
  token                   = "foo"
  property_mappings       = [data.authentik_property_mapping_scim.user.id]
  property_mappings_group = [data.authentik_property_mapping_scim.group.id]
}
