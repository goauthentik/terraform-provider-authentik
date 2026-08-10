package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStagePassword(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStagePassword(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_password.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStagePassword(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_password.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStagePassword(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_password" "name" {
  name              = "%[1]s"
  backends = ["authentik.core.auth.InbuiltBackend"]
}
`, name)
}
