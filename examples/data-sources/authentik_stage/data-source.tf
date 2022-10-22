# To get the ID of a stage by name

data "authentik_stage" "default-authentication-identification" {
  name = "default-authentication-identification"
}

# Then use `data.authentik_stage.default-authentication-identification.id`
