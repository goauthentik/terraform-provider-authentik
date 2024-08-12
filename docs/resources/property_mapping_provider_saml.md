---
page_title: "authentik_property_mapping_provider_saml Resource - terraform-provider-authentik"
subcategory: "Customization"
description: |-
  Manage SAML Provider Property mappings
---

# authentik_property_mapping_provider_saml (Resource)

Manage SAML Provider Property mappings

## Example Usage

```terraform
# Create a custom SAML provider property mapping

resource "authentik_property_mapping_provider_saml" "saml-aws-rolessessionname" {
  name       = "SAML AWS RoleSessionName"
  saml_name  = "https://aws.amazon.com/SAML/Attributes/RoleSessionName"
  expression = "return user.email"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `expression` (String)
- `name` (String)
- `saml_name` (String)

### Optional

- `friendly_name` (String)

### Read-Only

- `id` (String) The ID of this resource.