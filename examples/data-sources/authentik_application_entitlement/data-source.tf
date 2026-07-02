# To get the ID of an application entitlement by name

data "authentik_application" "example-application" {
  slug = "example-application"
}

data "authentik_application_entitlement" "example-entitlement" {
  app  = data.authentik_application.example-application.id
  name = "example-entitlement"
}

# Then use `data.authentik_application_entitlement.example-entitlement.id`
