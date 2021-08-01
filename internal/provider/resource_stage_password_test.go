package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceStagePassword(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStagePassword,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_password.name", "name", "test"),
				),
			},
		},
	})
}

const testAccResourceStagePassword = `
resource "authentik_stage_password" "name" {
  name              = "test"
  backends = ["django.contrib.auth.backends.ModelBackend"]
}
`
