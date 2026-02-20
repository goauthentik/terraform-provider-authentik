package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceStageRedirect(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageRedirect(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_redirect.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageRedirect(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_redirect.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageRedirect(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_redirect" "name" {
  name = "%[1]s"
  mode = "static"
  target_static = "https://goauthentik.io"
}
`, name)
}
