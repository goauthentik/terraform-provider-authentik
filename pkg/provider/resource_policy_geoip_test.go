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

func TestAccResourcePolicyGeoIP(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckPolicyGeoIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyGeoIP(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_geoip.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyGeoIP(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_geoip.name", "name", rName+"test"),
				),
			},
			// The config above leaves both check_* booleans unset, so they are null in state
			// and import round-trips them exactly. It also covers the five defaulted
			// numeric attributes, which would import as null under the pre-discovery-#5
			// pattern since the API can legitimately return 0 for them.
			{
				ResourceName:      "authentik_policy_geoip.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// These two steps together are the end-to-end regression test for the
			// GetP[bool] bug: switch check_history_distance on, then back off. Under SDKv2
			// the second step was impossible - GetP returned nil for false, omitempty
			// dropped the key, the API kept true, and Read returned true, so the plan after
			// apply was never empty. terraform-plugin-testing fails a step whose post-apply
			// plan is non-empty, so step 5 passing is the assertion that the check can
			// actually be disabled.
			{
				Config: testAccResourcePolicyGeoIPCheck(rName+"test", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_geoip.name", "check_history_distance", "true"),
				),
			},
			{
				Config: testAccResourcePolicyGeoIPCheck(rName+"test", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_geoip.name", "check_history_distance", "false"),
				),
			},
		},
	})
}

func testAccCheckPolicyGeoIPDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_policy_geoip" {
			continue
		}
		_, hr, err := c.PoliciesAPI.PoliciesGeoipRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("geoip policy %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourcePolicyGeoIP(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_geoip" "name" {
  name              = "%[1]s"
  asns = [123]
  countries = ["DE"]
}
`, name)
}

func testAccResourcePolicyGeoIPCheck(name string, checkHistoryDistance bool) string {
	return fmt.Sprintf(`
resource "authentik_policy_geoip" "name" {
  name      = "%[1]s"
  asns      = [123]
  countries = ["DE"]

  check_history_distance = %[2]t
}
`, name, checkHistoryDistance)
}
