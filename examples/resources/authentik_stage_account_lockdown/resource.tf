# Create an account lockdown stage

resource "authentik_stage_account_lockdown" "name" {
  name                  = "account-lockdown"
  deactivate_user       = true
  set_unusable_password = true
  delete_sessions       = true
  revoke_tokens         = true
}
