# Create a WebAuthn setup stage

resource "authentik_stage_authenticator_webauthn" "name" {
  name = "webauthn-setup"
}

# Create a WebAuthn setup which allows limited WebAuthn device types

data "authentik_webauthn_device_type" "yubikey" {
  description = "YubiKey 5C"
}

resource "authentik_stage_authenticator_webauthn" "name" {
  name = "webauthn-setup"
  device_type_restrictions = [
    data.authentik_webauthn_device_type.yubikey.aaguid,
  ]
}
