# To get the ID of a Telegram Source Property mapping

data "authentik_property_mapping_source_telegram" "test" {
  name = "custom-field"
}

# Then use `data.authentik_property_mapping_source_telegram.test.id`
