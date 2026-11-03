package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Step 1 leaves filter_group unset, which used to panic on read (#933).
func TestAccResourceProviderGoogleWorkspace(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderGoogleWorkspace(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_google_workspace.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_google_workspace.name", "filter_group", ""),
				),
			},
			{
				Config: testAccResourceProviderGoogleWorkspace(rName+"test", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_google_workspace.name", "name", rName+"test"),
					resource.TestCheckResourceAttrPair("authentik_provider_google_workspace.name", "filter_group", "authentik_group.group", "id"),
				),
			},
		},
	})
}

func testAccResourceProviderGoogleWorkspace(name string, filterGroup bool) string {
	fg := ""
	if filterGroup {
		fg = "filter_group = authentik_group.group.id"
	}
	return fmt.Sprintf(`
resource "authentik_group" "group" {
  name = "%[1]s"
}

resource "authentik_provider_google_workspace" "name" {
  name = "%[1]s"
  default_group_email_domain = "goauthentik.io"
  %[2]s
}
`, name, fg)
}
