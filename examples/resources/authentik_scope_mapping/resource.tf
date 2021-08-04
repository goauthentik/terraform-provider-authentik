# Create a scope mapping

resource "authentik_scope_mapping" "name" {
  name       = "minio"
  scope_name = "minio"
  expression = <<EOF
return {
  "policy": "readwrite",
}
EOF
}
