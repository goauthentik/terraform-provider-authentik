# Create a prompt field

resource "authentik_stage_prompt_field" "field" {
  field_key = "username"
  label     = "Username"
  type      = "username"
}
