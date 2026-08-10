package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageMutualTLS(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: providerFactories,
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
