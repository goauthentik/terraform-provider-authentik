# Create a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}

# Create a user that is member of a group

resource "authentik_group" "group" {
  name = "group-name"
}

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
  groups   = [authentik_group.group.id]
}
