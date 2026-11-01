resource "authentik_request_rule" "example" {
  name                       = "Example access request rule"
  policy_engine_mode         = "any"
  notification_mode          = "all"
  min_reviewers              = 1
  min_reviewers_is_per_group = false
}
