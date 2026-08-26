resource "authentik_request_rule" "example" {
  name = "Example access request rule"
}

resource "authentik_application" "primary" {
  name = "example app"
  slug = "example-app"
}

resource "authentik_application" "secondary" {
  name = "example app addon"
  slug = "example-app-addon"
}

resource "authentik_request_rule_binding" "example" {
  rule   = authentik_request_rule.example.id
  target = authentik_application.primary.uuid
}

# Make the addon requestable alongside the primary application, under the same rule
resource "authentik_request_rule_child_binding" "example" {
  binding = authentik_request_rule_binding.example.id
  target  = authentik_application.secondary.uuid
}
