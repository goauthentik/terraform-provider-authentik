# Create a policy binding for a resource

resource "authentik_policy_expression" "policy" {
  name       = "example"
  expression = "return True"
}

resource "authentik_application" "name" {
  name = "test app"
  slug = "test-app"
}

resource "authentik_policy_binding" "app-access" {
  target = authentik_application.name.uuid
  policy = authentik_policy_expression.policy.id
  order  = 0
}

# Create a binding to a group

data "authentik_group" "admins" {
  name = "authentik Admins"
}

resource "authentik_application" "name" {
  name = "test app"
  slug = "test-app"
}

resource "authentik_policy_binding" "app-access" {
  target = authentik_application.name.uuid
  group  = data.authentik_group.admins.id
  order  = 0
}
