# To get the ID of a scope mapping

data "authentik_property_mapping_provider_scope" "test" {
  # Search by name, by managed field or by scope_name
  # name    = "authentik default OAuth Mapping: Proxy outpost"
  managed = "goauthentik.io/providers/proxy/scope-proxy"
}

# Then use `data.authentik_property_mapping_provider_scope.test.id`

# Or, to get the IDs of multiple mappings

data "authentik_property_mapping_provider_scope" "test" {
  managed_list = [
    "goauthentik.io/providers/oauth2/scope-email",
    "goauthentik.io/providers/oauth2/scope-openid"
  ]
}

# Then use data.authentik_property_mapping_provider_scope.test.ids
