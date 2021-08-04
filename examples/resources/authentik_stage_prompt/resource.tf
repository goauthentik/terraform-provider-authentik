# Create a prompt stage with 1 field

resource "authentik_stage_prompt_field" "field" {
  field_key = "username"
  label     = "Username"
  type      = "username"
}
resource "authentik_stage_prompt" "name" {
  name = "test"
  fields = [
    resource.authentik_stage_prompt_field.field.id,
  ]
}
