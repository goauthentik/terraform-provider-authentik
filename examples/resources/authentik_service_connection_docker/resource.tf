# Create a local docker connection

resource "authentik_service_connection_docker" "local" {
  name  = "local"
  local = true
}

# Create a remote docker connection

resource "authentik_certificate_key_pair" "tls-auth" {
  name             = "docker-tls-auth"
  certificate_data = "..."
  key_data         = "..."
}

resource "authentik_certificate_key_pair" "tls-verification" {
  name             = "docker-tls-verification"
  certificate_data = "..."
}

resource "authentik_service_connection_docker" "remote-host" {
  name               = "remote-host"
  url                = "http://1.2.3.4:2368"
  tls_verification   = authentik_certificate_key_pair.tls-auth.id
  tls_authentication = authentik_certificate_key_pair.tls-verification.id
}
