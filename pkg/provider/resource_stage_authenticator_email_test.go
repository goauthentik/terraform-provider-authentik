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

func TestAccResourceStageAuthenticatorEmail(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageAuthenticatorEmailDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorEmail(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_email.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorEmail(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_email.name", "name", rName+"test"),
				),
			},
			// password is write-only server-side (H4) and cannot be imported; every other
			// attribute must round-trip exactly.
			{
				ResourceName:            "authentik_stage_authenticator_email.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckStageAuthenticatorEmailDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_authenticator_email" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesAuthenticatorEmailRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("email authenticator stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageAuthenticatorEmail(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_email" "name" {
  name              = "%[1]s"
}
`, name)
}
