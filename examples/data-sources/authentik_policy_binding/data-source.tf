# To get the ID of a policy binding by application and group

data "authentik_application" "example-application" {
  slug = "example-application"
}

data "authentik_group" "example-users" {
  name = "example-users"
}

data "authentik_policy_binding" "example-binding" {
  target = data.authentik_application.example-application.uuid
  group  = data.authentik_group.example-users.id
}

# To get the ID of a policy binding by application and policy

data "authentik_application" "example-application" {
  slug = "example-application"
}

data "authentik_policy_expression" "example-policy" {
  name = "example-policy"
}

data "authentik_policy_binding" "example-binding" {
  target = data.authentik_application.example-application.uuid
  policy = data.authentik_policy_expression.example-policy.id
}

# To get the ID of a policy binding by application and user

data "authentik_application" "example-application" {
  slug = "example-application"
}

data "authentik_user" "example-user" {
  username = "example-user"
}

data "authentik_policy_binding" "example-binding" {
  target = data.authentik_application.example-application.uuid
  policy = data.authentik_user.example-user.id
}

# To get the ID of a policy binding where multiple matching app/group bindings exist
# (or policy, or user)

data "authentik_application" "example-app" {
  slug = "example-application"
}

data "authentik_group" "example-users" {
  name = "example-users"
}

data "authentik_policy_binding" "example-binding" {
  target = data.authentik_application.example-application.uuid
  group  = data.authentik_group.example-users.id
  order  = 10
}

# Then use `data.authentik_policy_binding.example-binding.id`
