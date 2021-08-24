# Create identification stage with sources and showing a password field

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_oauth" "name" {
  name                = "test"
  slug                = "test"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow     = data.authentik_flow.default-authorization-flow.id

  provider_type   = "discord"
  consumer_key    = "foo"
  consumer_secret = "bar"
}

resource "authentik_stage_password" "name" {
  name     = "test-pass"
  backends = ["authentik.core.auth.InbuiltBackend"]
}

resource "authentik_stage_identification" "name" {
  name           = "test-ident"
  user_fields    = ["username"]
  sources        = [authentik_source_oauth.name.uuid]
  password_stage = authentik_stage_password.name.id
}
