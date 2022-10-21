# To get the complete users list

data "authentik_users" "all" {
}

# Then use `data.authentik_users.all.users`

# Or, to filter according to a specific field

data "authentik_users" "admins" {
  is_superuser = true
}

# Then use `data.authentik_users.admins.users`
