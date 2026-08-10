package sdkprovider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceEventTransport(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEventTransport(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_event_transport.transport", "name", rName),
				),
			},
			{
				Config: testAccResourceEventTransport(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_event_transport.transport", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceEventTransport(name string) string {
	return fmt.Sprintf(`
resource "authentik_event_transport" "transport" {
  name = "%[1]s"
  mode = "webhook"
  send_once = true
  webhook_url = "https://foo.bar"
}
`, name)
}
