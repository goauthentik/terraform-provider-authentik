# Create a TOTP Setup stage

resource "authentik_stage_authenticator_totp" "name" {
  name = "totp-setup"
}
