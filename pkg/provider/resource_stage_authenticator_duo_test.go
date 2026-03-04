package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceStageAuthenticatorDuo(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
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
