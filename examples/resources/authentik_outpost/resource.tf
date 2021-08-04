# Create an outpost with a proxy provider

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_proxy" "proxy" {
  name               = "proxy"
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  external_host      = "http://foo.bar.baz"
  internal_host      = "http://internal.local"
}

resource "authentik_outpost" "outpost" {
  name = "test-outpost"
  protocol_providers = [
    authentik_provider_proxy.proxy.id
  ]
}
