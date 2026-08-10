package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageConsent(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageConsent(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_consent.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageConsent(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_consent.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageConsent(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_consent" "name" {
  name              = "%[1]s"
}
`, name)
}
