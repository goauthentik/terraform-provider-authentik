# Create a SAML Provider

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_saml" "name" {
  name               = "test-app"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  acs_url            = "http://localhost"
}

resource "authentik_application" "name" {
  name              = "test-app"
  slug              = "test-app"
  protocol_provider = authentik_provider_saml.name.id
}
