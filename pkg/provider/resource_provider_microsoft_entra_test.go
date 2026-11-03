package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Step 1 leaves filter_group unset, which used to panic on read (#933).
func TestAccResourceProviderMicrosoftEntra(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderMicrosoftEntra(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_microsoft_entra.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_microsoft_entra.name", "filter_group", ""),
				),
			},
			{
				Config: testAccResourceProviderMicrosoftEntra(rName+"test", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_microsoft_entra.name", "name", rName+"test"),
					resource.TestCheckResourceAttrPair("authentik_provider_microsoft_entra.name", "filter_group", "authentik_group.group", "id"),
				),
			},
		},
	})
}

func testAccResourceProviderMicrosoftEntra(name string, filterGroup bool) string {
	fg := ""
	if filterGroup {
		fg = "filter_group = authentik_group.group.id"
	}
	return fmt.Sprintf(`
resource "authentik_group" "group" {
  name = "%[1]s"
}

resource "authentik_provider_microsoft_entra" "name" {
  name = "%[1]s"
  client_id     = "foo"
  client_secret = "foo"
  tenant_id     = "bar"
  %[2]s
}
`, name, fg)
}
