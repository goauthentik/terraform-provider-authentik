package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageAuthenticatorEmail(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorEmail(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_email.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorEmail(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_email.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageAuthenticatorEmail(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_email" "name" {
  name              = "%[1]s"
}
`, name)
}
