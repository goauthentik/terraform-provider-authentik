# Create a token for a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}

resource "authentik_token" "default" {
  identifier  = "my-token"
  user        = authentik_user.name.id
  description = "My secret token"
  expires     = "2025-01-01T15:04:05Z"
}
