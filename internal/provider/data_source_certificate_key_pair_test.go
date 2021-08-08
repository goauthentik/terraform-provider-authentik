package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCertificateKeyPair(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
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
