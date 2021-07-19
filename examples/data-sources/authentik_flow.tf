# To get the ID of a flow by slug

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

# Then use `data.authentik_flow.default-authorization-flow.id`
