provider "authentik" {
  url   = "https://authentik.company"
  token = "foo-bar"
  # Optionally set insecure to ignore TLS Certificates
  # insecure = true
  # Optionally add extra headers
  # headers {
  #   X-my-header = "foo"
  # }
}
