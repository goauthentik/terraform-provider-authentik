package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceWebAuthnDeviceType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWebAuthnDeviceType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_webauthn_device_type.op", "description", "1Password"),
					resource.TestCheckResourceAttr("data.authentik_webauthn_device_type.op", "aaguid", "bada5566-a7aa-401f-bd96-45619a55120d"),
				),
			},
		},
	})
}

const testAccDataSourceWebAuthnDeviceType = `
data "authentik_webauthn_device_type" "op" {
  description = "1Password"
}
`
