# Create a scope mapping

resource "authentik_property_mapping_provider_scope" "name" {
  name       = "minio"
  scope_name = "minio"
  expression = <<EOF
return {
  "policy": "readwrite",
}
EOF
}
