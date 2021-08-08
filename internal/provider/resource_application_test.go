package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceApplication(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationSimple(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "protocol_provider", "0"),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_launch_url", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_description", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_publisher", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "policy_engine_mode", string(api.POLICYENGINEMODE_ANY)),
				),
			},
			{
				Config: testAccResourceApplicationSimple(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "protocol_provider", "0"),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_launch_url", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_description", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_publisher", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "policy_engine_mode", string(api.POLICYENGINEMODE_ANY)),
				),
			},
			{
				Config:      testAccResourceApplicationSimple(rName + "test+"),
				ExpectError: regexp.MustCompile("consisting of letters, numbers, underscores or hyphens"),
			},
		},
	})
}

func testAccResourceApplicationSimple(name string) string {
	return fmt.Sprintf(`
resource "authentik_application" "name" {
  name              = "%[1]s"
  slug              = "%[1]s"
  meta_icon = "http://localhost/%[1]s"
}
`, name)
}
