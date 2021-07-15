terraform {
  required_providers {
    authentik = {
      version = "0.2"
      source  = "goauthentik/authentik"
    }
  }
}

provider "authentik" {
  url = "https://beryjuorg-dev.my.goauthentik.io"
}

resource "authentik_provider_proxy" "proxy" {
  name = "proxy"
  authorization_flow = "f75c8903-f82e-4783-a556-572134ac0bf1"
  external_host = "http://foo.bar.baz"
  internal_host = "http://internal.local"
}

resource "authentik_application" "name" {
  name = "test app"
  slug = "test-app"
  # protocol_provider = authentik_provider_proxy.proxy.id
}
