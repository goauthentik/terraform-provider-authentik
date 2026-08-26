resource "authentik_user" "example" {
  username = "leaving-user"
  name     = "Leaving User"
}

resource "authentik_user_offboarding" "example" {
  user            = authentik_user.example.id
  scheduled_at    = "2026-12-31T00:00:00Z"
  action          = "deactivate"
  revoke_sessions = true
  revoke_tokens   = true
}
