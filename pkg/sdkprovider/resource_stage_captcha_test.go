package sdkprovider_test

import (
	"context"
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceStageCaptcha(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageCaptchaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageCaptcha(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_captcha.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageCaptcha(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_captcha.name", "name", rName+"test"),
				),
			},
			{
				ResourceName:            "authentik_stage_captcha.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func testAccCheckStageCaptchaDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_captcha" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesCaptchaRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("stage_captcha %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageCaptcha(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_captcha" "name" {
  name              = "%[1]s"
  private_key = "foo"
  public_key = "bar"
}
`, name)
}
