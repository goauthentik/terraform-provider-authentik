# Create a local kubernetes connection

resource "authentik_service_connection_kubernetes" "local" {
  name  = "local"
  local = true
}

# Create a remote kubernetes connection

resource "authentik_service_connection_kubernetes" "remote-test-cluster" {
  name       = "test-cluster"
  kubeconfig = <<EOF
kind: Config
users: [...]
EOF
}
