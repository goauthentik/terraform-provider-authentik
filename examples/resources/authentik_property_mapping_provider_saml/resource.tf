# Create a custom SAML provider property mapping

resource "authentik_property_mapping_provider_saml" "saml-aws-rolessessionname" {
  name       = "SAML AWS RoleSessionName"
  saml_name  = "https://aws.amazon.com/SAML/Attributes/RoleSessionName"
  expression = "return user.email"
}
