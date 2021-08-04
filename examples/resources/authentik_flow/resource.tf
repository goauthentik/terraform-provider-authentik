# Create a flow with a stage attached

resource "authentik_stage_dummy" "name" {
  name = "test-stage"
}

resource "authentik_flow" "flow" {
  name        = "test-flow"
  title       = "Test flow"
  slug        = "test-flow"
  designation = "authorization"
}

resource "authentik_flow_stage_binding" "dummy-flow" {
  target = authentik_flow.flow.uuid
  stage  = authentik_stage_dummy.name.id
  order  = 0
}
