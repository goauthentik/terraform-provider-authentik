package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageRedirect(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
