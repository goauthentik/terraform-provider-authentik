# Create a super-user group with a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}
resource "authentik_group" "group" {
  name         = "tf_admins"
  users        = [authentik_user.name.id]
  is_superuser = true
}
