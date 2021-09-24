# Create a custom Notification transport mapping

resource "authentik_property_mapping_notification" "name" {
  name       = "custom-field"
  expression = "return {\"foo\": context['foo']}"
}
