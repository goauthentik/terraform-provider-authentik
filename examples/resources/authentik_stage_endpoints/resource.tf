resource "authentik_endpoints_connector_agent" "name" {
  name = "agent"
}

resource "authentik_stage_endpoints" "name" {
  name      = "agent-connector"
  connector = authentik_endpoints_connector_agent.name.id
}
