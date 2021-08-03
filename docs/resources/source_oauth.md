---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "authentik_source_oauth Resource - terraform-provider-authentik"
subcategory: ""
description: |-
  
---

# authentik_source_oauth (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **authentication_flow** (String)
- **consumer_key** (String)
- **consumer_secret** (String, Sensitive)
- **enrollment_flow** (String)
- **name** (String)
- **provider_type** (String)
- **slug** (String)

### Optional

- **access_token_url** (String)
- **authorization_url** (String)
- **enabled** (Boolean) Defaults to `true`.
- **id** (String) The ID of this resource.
- **policy_engine_mode** (String) Defaults to `any`.
- **profile_url** (String)
- **request_token_url** (String)
- **user_matching_mode** (String) Defaults to `identifier`.
- **uuid** (String)

