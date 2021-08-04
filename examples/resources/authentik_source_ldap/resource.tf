# Create LDAP Source

resource "authentik_source_ldap" "name" {
  name = "ldap-test"
  slug = "ldap-test"

  server_uri    = "ldaps://1.2.3.4"
  bind_cn       = "foo"
  bind_password = "bar"
  base_dn       = "dn=foo"
}
