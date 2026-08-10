package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageAccountLockdown(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAccountLockdown(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "deactivate_user", "true"),
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "revoke_tokens", "true"),
				),
			},
			{
				Config: testAccResourceStageAccountLockdown(rName+"test", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "deactivate_user", "false"),
					resource.TestCheckResourceAttr("authentik_stage_account_lockdown.name", "revoke_tokens", "false"),
				),
			},
		},
	})
}

func testAccResourceStageAccountLockdown(name string, lockdown bool) string {
	return fmt.Sprintf(`
resource "authentik_stage_account_lockdown" "name" {
  name                  = "%[1]s"
  deactivate_user       = %[2]t
  set_unusable_password = %[2]t
  delete_sessions       = %[2]t
  revoke_tokens         = %[2]t
}
`, name, lockdown)
}
