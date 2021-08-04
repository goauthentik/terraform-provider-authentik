# Create expression policys

resource "authentik_policy_expression" "name" {
  name       = "expression"
  expression = "return True"
}
