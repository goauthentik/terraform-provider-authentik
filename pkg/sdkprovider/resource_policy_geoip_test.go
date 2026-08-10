package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyGeoIP(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
		},
	})
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
