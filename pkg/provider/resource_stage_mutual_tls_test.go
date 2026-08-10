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

func TestAccResourceStageMutualTLS(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckStageMutualTLSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageMutualTLS(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_mutual_tls.name", "name", rName),
				),
			},
			{
				Config: testAccResourceStageMutualTLS(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_mutual_tls.name", "name", rName+"test"),
				),
			},
			// certificate_authorities is Optional-only and unset in the config above,
			// so this also pins H1's list variant: the attribute must import back as
			// null, not as an empty list.
			{
				ResourceName:      "authentik_stage_mutual_tls.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckStageMutualTLSDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_stage_mutual_tls" {
			continue
		}
		_, hr, err := c.StagesAPI.StagesMtlsRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("mutual TLS stage %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceStageMutualTLS(name string) string {
	return fmt.Sprintf(`
resource "authentik_stage_mutual_tls" "name" {
  name = "%[1]s"
}
`, name)
}
