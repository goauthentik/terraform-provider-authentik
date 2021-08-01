package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageIdentification(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageIdentification(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_identification.name", "name", rName),
				),
			},
		},
	})
}

func testAccResourceStageIdentification(name string) string {
	return fmt.Sprintf(`
# TODO: Source and sources field
resource "authentik_stage_identification" "name" {
  name              = "%s"
  user_fields = ["username"]
}
`, name)
}
