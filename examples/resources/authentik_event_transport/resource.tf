# Create a notification transport to a generic webhook URL
resource "authentik_event_transport" "transport" {
  name        = "my-transport"
  mode        = "webhook"
  send_once   = true
  webhook_url = "https://foo.bar"
}

# Create a notification transport to slack/discord
resource "authentik_event_transport" "transport" {
  name        = "my-transport"
  mode        = "webhook_slack"
  send_once   = true
  webhook_url = "https://discord.com/...."
}
