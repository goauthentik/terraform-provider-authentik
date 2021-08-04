# Create a static TOTP Setup stage

resource "authentik_stage_authenticator_static" "name" {
  name = "static-totp-setup"
}
