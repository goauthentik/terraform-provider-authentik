data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name              = "authentik Self-signed Certificate"
  fetch_key         = false
  fetch_certificate = false
}

resource "authentik_provider_oauth2" "example" {
  name               = "dcr-example"
  client_id          = "dcr-example"
  signing_key        = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow  = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_provider_oauth2_dcr" "example" {
  oauth2_provider        = authentik_provider_oauth2.example.id
  access_token_validity  = "minutes=5"
  refresh_token_validity = "days=30"
  allowed_grant_types    = ["authorization_code", "refresh_token"]
}
