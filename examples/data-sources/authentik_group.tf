# To get the ID of a group by name

data "authentik_group" "admins" {
  name = "authentik Admins"
}

# Then use `data.authentik_group.admins.id`
