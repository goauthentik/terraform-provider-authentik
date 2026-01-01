# Modify the schedule of a SCIM provider

resource "authentik_provider_scim" "name" {
  name  = "name"
  url   = "http://localhost"
  token = "foo"
}

resource "authentik_task_schedule" "default" {
  app_model = "authentik_providers_scim.scimprovider"
  model_id  = authentik_provider_scim.name.id

  crontab = "6 */4 * * 2"
}
