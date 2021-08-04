# Create expiry policy

resource "authentik_policy_expiry" "name" {
  name = "expiry"
  days = 3
}
