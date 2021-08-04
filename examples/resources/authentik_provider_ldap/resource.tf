# Create an LDAP Provider

data "authentik_flow" "default-authentication-flow" {
  slug = "default-authentication-flow"
}

resource "authentik_provider_ldap" "name" {
  name      = "ldap-app"
  base_dn   = "dc=ldap,dc=goauthentik,dc=io"
  bind_flow = data.authentik_flow.default-authentication-flow.id
}

resource "authentik_application" "name" {
  name              = "ldap-app"
  slug              = "ldap-app"
  protocol_provider = authentik_provider_ldap.name.id
}
