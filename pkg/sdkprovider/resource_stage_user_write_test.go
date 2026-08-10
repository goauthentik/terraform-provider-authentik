package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageUserWrite(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserWrite(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "create_users_as_inactive", "false"),
				),
			},
			{
				Config: testAccResourceStageUserWrite(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_stage_user_write.name", "create_users_as_inactive", "false"),
				),
			},
		},
	})
}

func testAccResourceStageUserWrite(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_write" "name" {
  name              = "%[1]s"
  create_users_as_inactive = false
}
`, name)
}
