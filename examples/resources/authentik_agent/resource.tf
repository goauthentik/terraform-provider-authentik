resource "authentik_agent" "example" {
  username        = "example-agent"
  name            = "Example Agent"
  policy_behavior = "mirror"
  is_active       = true
}
