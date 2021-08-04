# Create email stage for email verification, uses global settings by default

resource "authentik_stage_email" "name" {
  name = "email-verification"
}
