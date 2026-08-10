package sdkprovider_test

import (
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceStage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStageSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_stage.default-authentication-identification", "name", "default-authentication-identification"),
				),
			},
		},
	})
}

const testAccDataSourceStageSimple = `
data "authentik_stage" "default-authentication-identification" {
  name = "default-authentication-identification"
}
`
