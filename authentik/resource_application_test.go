package authentik

import (
	"testing"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceApplication(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplicationSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_application.name", "name", "acc-test-app"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", "acc-test-app"),
					resource.TestCheckResourceAttr("authentik_application.name", "protocol_provider", "0"),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_launch_url", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_description", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "meta_publisher", ""),
					resource.TestCheckResourceAttr("authentik_application.name", "policy_engine_mode", string(api.POLICYENGINEMODE_ANY)),
				),
			},
		},
	})
}

const testAccResourceApplicationSimple = `
resource "authentik_application" "name" {
  name              = "acc-test-app"
  slug              = "acc-test-app"
}
`
