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

func TestAccResourceStagePassword(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStagePasswordDestroy,
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
			// The H4/H1 detector for this batch: it asserts Read reproduces
			// post-apply state exactly, which is what catches both an attribute
			// fromAPI forgets to write and a defaulted attribute wrongly mapped to
			// null. allow_show_password (Default false) and
			// failed_attempts_before_cancel (Default 5) are unset in the config
			// above, so they only round-trip correctly because fromAPI writes them
			// verbatim rather than through the prior-aware helpers.
			{
				ResourceName:      "authentik_stage_password.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckStagePasswordDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_password" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesPasswordRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("password stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStagePassword(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_password" "name" {
  name              = "%[1]s"
  backends = ["authentik.core.auth.InbuiltBackend"]
}
`, name)
}
