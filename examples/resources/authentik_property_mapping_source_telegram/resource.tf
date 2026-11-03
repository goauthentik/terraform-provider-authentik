# Create a custom Telegram source property mapping

resource "authentik_property_mapping_source_telegram" "name" {
  name       = "custom-field"
  expression = "return {\"username\": data}"
}
