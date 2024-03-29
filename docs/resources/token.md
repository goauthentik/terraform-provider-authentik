---
page_title: "authentik_token Resource - terraform-provider-authentik"
subcategory: "Directory"
description: |-
  
---

# authentik_token (Resource)



## Example Usage

```terraform
# Create a token for a user

resource "authentik_user" "name" {
  username = "user"
  name     = "User"
}

resource "authentik_token" "default" {
  identifier  = "my-token"
  user        = authentik_user.name.id
  description = "My secret token"
  expires     = "2025-01-01T15:04:05Z"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `identifier` (String)
- `user` (Number)

### Optional

- `description` (String)
- `expires` (String)
- `expiring` (Boolean) Defaults to `true`.
- `intent` (String) Allowed values:
  - `verification`
  - `api`
  - `recovery`
  - `app_password`
 Defaults to `api`.
- `retrieve_key` (Boolean) Defaults to `false`.

### Read-Only

- `expires_in` (Number) Generated.
- `id` (String) The ID of this resource.
- `key` (String, Sensitive) Generated.
