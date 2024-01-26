# Create/manage a default brand

resource "authentik_brand" "default" {
  domain         = "."
  default        = true
  branding_title = "test"
}
