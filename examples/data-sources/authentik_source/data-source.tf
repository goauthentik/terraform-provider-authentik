data "authentik_source" "inbuilt" {
  managed = "goauthentik.io/sources/inbuilt"
}

# Then use `data.authentik_source.inbuilt.uuid`
