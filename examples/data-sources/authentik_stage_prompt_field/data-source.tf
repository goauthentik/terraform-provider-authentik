# To get the ID of a prompt field by name

data "authentik_stage_prompt_field" "default_user_settings_field_email" {
  name = "default-user-settings-field-email"
}

# Then use `data.authentik_stage_prompt_field.default_user_settings_field_email.id`
