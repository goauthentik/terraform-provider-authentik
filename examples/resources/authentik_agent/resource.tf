resource "authentik_user" "example_agent_parent" {
  username = "example-agent-parent"
  name     = "Example Agent Parent"
}

resource "authentik_agent" "example" {
  username        = "example-agent"
  name            = "Example Agent"
  parent          = authentik_user.example_agent_parent.id
  policy_behavior = "mirror"
  is_active       = true
}
