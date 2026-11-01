resource "authentik_object_attribute" "example" {
  object_type = "authentik_core.user"
  key         = "employee_id"
  label       = "Employee ID"
  type        = "text"
  is_unique   = true
  is_required = false
}
