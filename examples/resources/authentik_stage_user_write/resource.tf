# Create a user write stage

resource "authentik_stage_user_write" "name" {
  name                     = "user-write"
  create_users_as_inactive = false
}
