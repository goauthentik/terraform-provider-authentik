# To get the ID of a Kubernetes Service Connection by name

data "authentik_service_connection_kubernetes" "local" {
  name = "Local Kubernetes Cluster"
}

# Then use `data.authentik_service_connection_kubernetes.local.id`
