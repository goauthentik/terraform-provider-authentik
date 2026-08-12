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

func TestAccResourceStageEmail(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageEmail(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_email.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageEmail(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_email.name", "name", rName+"test"),
				),
			},
			// This resource has the most defaulted attributes in the batch (eleven), so
			// it is the widest check that fromAPI maps every one of them verbatim rather
			// than through the prior-aware helpers - all eleven are unset in the config
			// above and so would import as null under the old pattern. password is
			// write-only server-side (H4) and cannot be imported.
			{
				ResourceName:            "authentik_stage_email.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckStageEmailDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_email" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesEmailRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("email stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageEmail(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_email" "name" {
  name              = "%[1]s"
}
`, name)
}
