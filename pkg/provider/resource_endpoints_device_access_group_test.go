package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestAccResourceEndpointsDeviceAccessGroup(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEndpointsDeviceAccessGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_device_access_group.group", "name", rName),
				),
			},
			{
				Config: testAccResourceEndpointsDeviceAccessGroup(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_endpoints_device_access_group.group", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceEndpointsDeviceAccessGroup(name string) string {
	return fmt.Sprintf(`
resource "authentik_endpoints_device_access_group" "group" {
  name = "%[1]s"
}
`, name)
}
