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

func TestAccResourceStageAuthenticatorDuo(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageAuthenticatorDuoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageAuthenticatorDuo(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_duo.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageAuthenticatorDuo(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_authenticator_duo.name", "name", rName+"test"),
				),
			},
			// The strongest H4 probe in this batch: the config above actually sets
			// client_secret, and it is Required, so if fromAPI ever stopped relying on
			// the model being seeded from prior state the attribute would be nulled out
			// of state entirely and the apply would fail. Both secrets are write-only
			// server-side, so import genuinely cannot recover them - exactly as under
			// SDKv2 - hence the ignore list; every other attribute must round-trip.
			{
				ResourceName:            "authentik_stage_authenticator_duo.name",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret", "admin_secret_key"},
			},
		},
	})
}

func testAccCheckStageAuthenticatorDuoDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_authenticator_duo" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesAuthenticatorDuoRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("duo authenticator stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageAuthenticatorDuo(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_authenticator_duo" "name" {
  name              = "%[1]s"
  client_id = "foo"
  client_secret = "bar"
  api_hostname = "http://foo.bar.baz"
}
`, name)
}
