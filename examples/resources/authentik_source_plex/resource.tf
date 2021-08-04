# Create a plex source

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_plex" "name" {
  name                = "plex"
  slug                = "plex"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow     = data.authentik_flow.default-authorization-flow.id
  client_id           = "foo-bar-baz"
  plex_token          = "foo"
}
