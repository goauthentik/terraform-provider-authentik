package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageAuthenticatorStatic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorStatic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_static.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorStatic(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_static.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorStatic(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_static" "name" {
  name              = "%[1]s"
}
`, name)
}
