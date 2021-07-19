# Create a policy binding for a resource

resource "authentik_application" "name" {
  name = "test app"
  slug = "test-app"
}

resource "authentik_policy_binding" "app-access" {
  target = authentik_application.name.id
  group  = data.authentik_group.admins.id
  order  = 0
}
