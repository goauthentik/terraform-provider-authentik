# Create a SAML Source

data "authentik_flow" "default-source-pre-authentication" {
  slug = "default-source-pre-authentication"
}

data "authentik_flow" "default-source-authentication" {
  slug = "default-source-authentication"
}

data "authentik_flow" "default-source-enrollment" {
  slug = "default-source-enrollment"
}

resource "authentik_source_saml" "name" {
  name                    = "test-source"
  slug                    = "test-source"
  authentication_flow     = data.authentik_flow.default-source-authentication.id
  enrollment_flow         = data.authentik_flow.default-source-enrollment.id
  pre_authentication_flow = data.authentik_flow.default-source-pre-authentication.id
  sso_url                 = "http://localhost"
}
