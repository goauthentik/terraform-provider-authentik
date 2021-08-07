# To get the ID of a scope mapping

data "authentik_scope_mapping" "test" {
  # Search by name, by managed field or by scope_name
  # name    = "authentik default OAuth Mapping: Proxy outpost"
  managed = "goauthentik.io/providers/proxy/scope-proxy"
}

# Then use `data.authentik_scope_mapping.test.id`
