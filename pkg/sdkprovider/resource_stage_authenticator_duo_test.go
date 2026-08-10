package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageAuthenticatorDuo(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorDuo(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_duo.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorDuo(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_duo.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorDuo(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_duo" "name" {
  name              = "%[1]s"
  client_id = "foo"
  client_secret = "bar"
  api_hostname = "http://foo.bar.baz"
}
`, name)
}
