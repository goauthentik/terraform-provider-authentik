# Create/manage a default tenant

resource "authentik_tenant" "default" {
  domain         = "."
  default        = true
  branding_title = "test"
}
