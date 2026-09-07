package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// The Google Workspace and Microsoft Entra providers both expose a nullable
// filter_group attribute. Reading either back used to panic when filter_group
// was unset, so this exercises create + refresh for both with filter_group
// omitted. See https://github.com/goauthentik/terraform-provider-authentik/issues/933.
func TestAccResourceProviderDirectorySync(t *testing.T) {
	for _, tc := range []struct {
		name     string
		resource string
		config   func(name string) string
	}{
		{
			name:     "google_workspace",
			resource: "authentik_provider_google_workspace.name",
			config: func(name string) string {
				return fmt.Sprintf(`
resource "authentik_provider_google_workspace" "name" {
  name                       = "%[1]s"
  default_group_email_domain = "goauthentik.io"
}
`, name)
			},
		},
		{
			name:     "microsoft_entra",
			resource: "authentik_provider_microsoft_entra.name",
			config: func(name string) string {
				return fmt.Sprintf(`
resource "authentik_provider_microsoft_entra" "name" {
  name      = "%[1]s"
  client_id = "foo"
  tenant_id = "bar"
}
`, name)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
			resource.UnitTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				Steps: []resource.TestStep{
					{
						Config: tc.config(rName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(tc.resource, "name", rName),
						),
					},
					{
						Config: tc.config(rName + "test"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(tc.resource, "name", rName+"test"),
						),
					},
				},
			})
		})
	}
}
