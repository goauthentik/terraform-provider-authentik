package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageUserLogout(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogout(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_logout.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageUserLogout(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_logout.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceStageUserLogout(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_logout" "name" {
  name              = "%[1]s"
}
`, name)
}
