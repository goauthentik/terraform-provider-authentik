# Create an Authenticator validations tage

resource "authentik_stage_authenticator_validate" "name" {
  name                  = "authenticator-validate"
  device_classes        = ["static"]
  not_configured_action = "skip"
}
