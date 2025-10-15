# Create a SCIM Source

resource "authentik_source_scim" "name" {
  name = "test-source"
  slug = "test-source"
}
