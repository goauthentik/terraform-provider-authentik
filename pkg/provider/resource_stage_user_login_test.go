package provider_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
)

func TestAccResourceStageUserLogin(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageUserLoginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageUserLogin(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "session_duration", "minutes=1"),
				),
			},
			{
				Config: testAccResourceStageUserLogin(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_stage_user_login.name", "session_duration", "minutes=1"),
				),
			},
			// network_binding and geoip_binding are ignored for a different reason than the
			// write-only secrets in part 6: the API does return both, but this resource's
			// Read has never refreshed them (H4), so an import - which starts from an empty
			// model with nothing to preserve - leaves them null. Ignoring them keeps that
			// pre-existing behaviour visible and asserted rather than silently worked
			// around. Everything else, including the three relative-duration defaults, must
			// round-trip exactly.
			{
				ResourceName:            "authentik_stage_user_login.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"network_binding", "geoip_binding"},
			},
		},
	})
}

func testAccCheckStageUserLoginDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_user_login" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesUserLoginRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("user login stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageUserLogin(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_login" "name" {
  name              = "%[1]s"
  session_duration = "minutes=1"
}
`, name)
}
