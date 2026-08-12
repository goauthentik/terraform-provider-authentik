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

func TestAccResourceStageUserWrite(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageUserWriteDestroy,
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
			// A direct probe for discovery #5, and the reason this resource is the one
			// carrying the import step in this batch: create_users_as_inactive is set to
			// false against a Default of true, and user_path_template defaults to "",
			// which is also what the API returns for it. Both are values the prior-aware
			// helpers would import as null, so this step fails if fromAPI ever regresses
			// to using them for defaulted attributes.
			{
				ResourceName:      "authentik_stage_user_write.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckStageUserWriteDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_user_write" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesUserWriteRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("user write stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageUserWrite(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_user_write" "name" {
  name              = "%[1]s"
  create_users_as_inactive = false
}
`, name)
}
