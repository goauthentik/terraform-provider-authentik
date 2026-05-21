# Create a Google Chrome endpoint connector

resource "authentik_endpoints_google_chrome_connector" "name" {
  name        = "google-chrome"
  enabled     = true
  credentials = file("${path.module}/service-account.json")
}
