# To get the name of a user by username

data "authentik_user" "akadmin" {
  username = "akadmin"
}

# Then use `data.authentik_group.akadmin.name`
