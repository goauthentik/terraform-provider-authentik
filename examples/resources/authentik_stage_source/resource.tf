# Create OAuth Source using an existing provider

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_oauth" "name" {
  name                = "discord"
  slug                = "discord"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow     = data.authentik_flow.default-authorization-flow.id

  provider_type   = "discord"
  consumer_key    = "foo"
  consumer_secret = "bar"
}

# Create a source stage using the source defined above
resource "authentik_stage_source" "name" {
  name   = "source-stage"
  source = authentik_source_oauth.name.id
}
