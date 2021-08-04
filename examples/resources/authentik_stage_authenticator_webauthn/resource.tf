# Create WebAuthn setup stage

resource "authentik_stage_authenticator_webauthn" "name" {
  name = "webauthn-setup"
}
