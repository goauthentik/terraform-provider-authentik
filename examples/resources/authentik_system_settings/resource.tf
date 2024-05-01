# Configure system settings

resource "authentik_system_settings" "settings" {
  default_user_change_username = true
}
