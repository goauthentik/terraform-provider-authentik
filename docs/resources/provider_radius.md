---
page_title: "authentik_provider_radius Resource - terraform-provider-authentik"
subcategory: "Applications"
description: |-
  
---

# authentik_provider_radius (Resource)



## Example Usage

```terraform
# Create a Radius Provider

data "authentik_flow" "default-authentication-flow" {
  slug = "default-authentication-flow"
}

resource "authentik_provider_radius" "name" {
  name               = "radius-app"
  authorization_flow = data.authentik_flow.default-authentication-flow.id
  client_networks    = "10.10.0.0/24"
  shared_secret      = "my-shared-secret"
}

resource "authentik_application" "name" {
  name              = "radius-app"
  slug              = "radius-app"
  protocol_provider = authentik_provider_radius.name.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `authorization_flow` (String)
- `name` (String)
- `shared_secret` (String, Sensitive)

### Optional

- `client_networks` (String) Defaults to `0.0.0.0/0, ::/0`.

### Read-Only

- `id` (String) The ID of this resource.

