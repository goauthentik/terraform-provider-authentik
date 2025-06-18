//go:build enterprise
// +build enterprise

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStageMutualTLS(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageMutualTLS(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_mutual_tls.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageMutualTLS(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_mutual_tls.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageMutualTLS(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_mutual_tls" "name" {
  name = "%[1]s"
}
`, name)
}
