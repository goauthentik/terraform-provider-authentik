# Create an Agent connector
resource "authentik_endpoints_connector_agent" "agent" {
  name = "agent"
}

resource "authentik_endpoints_connector_agent_enrollment_token" "token" {
  connector    = authentik_endpoints_connector_agent.agent.id
  name         = "my-enrollment token"
  expiring     = false
  retrieve_key = true
}

# then use the enrollment token via `authentik_endpoints_connector_agent_enrollment_token.token.key`
