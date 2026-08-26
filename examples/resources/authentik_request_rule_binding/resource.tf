resource "authentik_request_rule" "example" {
  name = "Example access request rule"
}

resource "authentik_application" "example" {
  name = "example app"
  slug = "example-app"
}

resource "authentik_request_rule_binding" "example" {
  rule               = authentik_request_rule.example.id
  target             = authentik_application.example.uuid
  expiry_pending     = "hours=8"
  expiry_granted_max = "hours=720"
}
