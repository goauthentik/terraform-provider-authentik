---
page_title: "authentik_scope_mapping Resource - terraform-provider-authentik"
subcategory: "Customization"
description: |-
  Manage Scope Provider Property mappings
  ~> This resource is deprecated. Migrate to authentik_property_mapping_provider_scope.
---

# authentik_scope_mapping (Resource)

Manage Scope Provider Property mappings

~> This resource is deprecated. Migrate to `authentik_property_mapping_provider_scope`.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `expression` (String)
- `name` (String)
- `scope_name` (String)

### Optional

- `description` (String)

### Read-Only

- `id` (String) The ID of this resource.
