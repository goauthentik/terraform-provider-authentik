# Create a proxy provider

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_proxy" "name" {
  name               = "test-app"
  internal_host      = "http://foo.bar.baz"
  external_host      = "http://internal.service"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
}

resource "authentik_application" "name" {
  name              = "test-app"
  slug              = "test-app"
  protocol_provider = authentik_provider_proxy.name.id
}
