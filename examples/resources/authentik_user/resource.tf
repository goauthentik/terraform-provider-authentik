# Create a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}
