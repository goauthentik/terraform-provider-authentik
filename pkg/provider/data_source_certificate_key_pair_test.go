package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccDataSourceCertificateKeyPair(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCertificateKeyPairSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_certificate_key_pair.generated", "name", "authentik Self-signed Certificate"),
					resource.TestCheckResourceAttr("data.authentik_certificate_key_pair.generated", "subject", "OU=Self-signed,O=authentik,CN=authentik Self-signed Certificate"),
				),
			},
		},
	})
}

const testAccDataSourceCertificateKeyPairSimple = `
data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
}
`
