# Create an application entitlement bound to a group
resource "authentik_application" "name" {
  name = "example-app"
  slug = "example-app"
}

resource "authentik_application_entitlement" "ent" {
  name        = "test-ent"
  application = authentik_application.name.id
}

resource "authentik_group" "group" {
  name = "test-ent-group"
}

resource "authentik_policy_binding" "test-ent-access" {
  target = authentik_application.name.uuid
  group  = authentik_group.group.id
  order  = 0
}
