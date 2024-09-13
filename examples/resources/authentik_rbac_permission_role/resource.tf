# Assign a global permission to a role

resource "authentik_rbac_role" "my-role" {
  name = "my-role"
}

resource "authentik_rbac_permission_role" "global-permission" {
  role       = authentik_rbac_role.my-role.id
  permission = "authentik_flows.inspect_flow"
}

# Assign an object permission to a role

resource "authentik_flow" "flow" {
  name        = "test-flow"
  title       = "Test flow"
  slug        = "test-flow"
  designation = "authorization"
}

resource "authentik_rbac_role" "my-role" {
  name = "my-role"
}

resource "authentik_rbac_permission_role" "global-permission" {
  role       = authentik_rbac_role.my-role.id
  model      = "authentik_flows.flow"
  permission = "inspect_flow"
  object_id  = authentik_flow.flow.uuid
}
