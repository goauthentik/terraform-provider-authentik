# To get the ID of a policy expression by name

data "authentik_policy_expression" "default_user_settings_authorization" {
  name = "default-user-settings-authorization"
}

# Then use `data.authentik_policy_expression.default_user_settings_authorization.id`
