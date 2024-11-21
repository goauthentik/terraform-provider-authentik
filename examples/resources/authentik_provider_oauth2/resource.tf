# Create an OAuth2 Provider

resource "authentik_provider_oauth2" "name" {
  name      = "grafana"
  client_id = "grafana"
  allowed_redirect_uris = [
    {
      matching_mode = "strict",
      url           = "http://localhost",
    }
  ]
}

resource "authentik_application" "name" {
  name              = "test app"
  slug              = "test-app"
  protocol_provider = authentik_provider_oauth2.name.id
}
