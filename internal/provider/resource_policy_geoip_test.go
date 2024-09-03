package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyGeoIP(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
