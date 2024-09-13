# Assign a global permission to a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}

resource "authentik_rbac_permission_user" "global-permission" {
  user       = authentik_user.name.id
  permission = "authentik_flows.inspect_flow"
}

# Assign an object permission to a user

resource "authentik_flow" "flow" {
  name        = "test-flow"
  title       = "Test flow"
  slug        = "test-flow"
  designation = "authorization"
}

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}

resource "authentik_rbac_permission_user" "global-permission" {
  user       = authentik_user.name.id
  model      = "authentik_flows.flow"
  permission = "inspect_flow"
  object_id  = authentik_flow.flow.uuid
}
