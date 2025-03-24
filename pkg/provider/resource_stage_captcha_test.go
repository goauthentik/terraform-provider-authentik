package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageCaptcha(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageCaptcha(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_captcha.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageCaptcha(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_captcha.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageCaptcha(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_captcha" "name" {
  name              = "%[1]s"
  private_key = "foo"
  public_key = "bar"
}
`, name)
}
