# To get the details of a brand by domain

data "authentik_brand" "authentik-default" {
  domain = "authentik-default"
}

# Then use `data.authentik_brand.authentik-default.domain`, `data.authentik_brand.authentik-default.branding_title`,
# `data.authentik_brand.authentik-default.branding_logo`, ...
