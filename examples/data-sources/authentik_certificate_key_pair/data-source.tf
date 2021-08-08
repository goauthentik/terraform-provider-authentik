# To get the the ID and other info about a certificate

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
}

# Then use `data.authentik_certificate_key_pair.generated.id`
