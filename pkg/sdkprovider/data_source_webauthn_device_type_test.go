package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceWebAuthnDeviceType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
