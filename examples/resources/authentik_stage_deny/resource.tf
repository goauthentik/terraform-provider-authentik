# Create deny stage, can be used with policies

resource "authentik_stage_deny" "name" {
  name = "deny"
}
