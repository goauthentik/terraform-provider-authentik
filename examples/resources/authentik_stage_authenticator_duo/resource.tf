# Create a duo setup stage

resource "authentik_stage_authenticator_duo" "name" {
  name          = "duo-setup"
  client_id     = "foo"
  client_secret = "bar"
  api_hostname  = "http://foo.bar.baz"
}
