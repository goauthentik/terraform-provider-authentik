# Create policy to match events

resource "authentik_policy_event_matcher" "name" {
  name      = "login-from-1.2.3.4"
  action    = "login"
  app       = "authentik.events"
  client_ip = "1.2.3.4"
}
