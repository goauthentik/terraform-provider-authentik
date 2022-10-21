# To get the complete groups list

data "authentik_groups" "all" {
}

# Then use `data.authentik_groups.all.groups`

# Or, to filter according to a specific field

data "authentik_groups" "admins" {
  is_superuser = true
}

# Then use `data.authentik_groups.admins.groups`
