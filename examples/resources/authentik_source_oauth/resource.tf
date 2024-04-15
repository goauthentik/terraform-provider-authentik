# Create OAuth Source using an existing provider

data "authentik_flow" "default-source-authentication" {
  slug = "default-source-authentication"
}
data "authentik_flow" "default-source-enrollment" {
  slug = "default-source-enrollment"
}

resource "authentik_source_oauth" "name" {
  name                = "discord"
  slug                = "discord"
  authentication_flow = data.authentik_flow.default-source-authentication.id
  enrollment_flow     = data.authentik_flow.default-source-enrollment.id

  provider_type   = "discord"
  consumer_key    = "foo"
  consumer_secret = "bar"
}
