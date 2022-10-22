# To get the details of a tenant by domain

data "authentik_tenant" "authentik-default" {
  domain = "authentik-default"
}

# Then use `data.authentik_tenant.authentik-default.domain`, `data.authentik_tenant.authentik-default.branding_title`,
# `data.authentik_tenant.authentik-default.branding_logo`, ...
