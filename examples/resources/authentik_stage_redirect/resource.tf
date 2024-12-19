# Create a static redirect stage

resource "authentik_stage_redirect" "static" {
  name          = "static-redirect"
  mode          = "static"
  target_static = "https://goauthentik.io"
}

# Create a flow redirect stage

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_stage_redirect" "flow" {
  name        = "flow-redirect"
  mode        = "flow"
  target_flow = data.authentik_flow.default-authorization-flow.id
}
