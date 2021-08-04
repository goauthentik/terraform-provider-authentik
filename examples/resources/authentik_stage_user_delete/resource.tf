# Create a user deletion stage

resource "authentik_stage_user_delete" "name" {
  name = "user-delete"
}
