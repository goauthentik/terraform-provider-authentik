package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTaskSchedule(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTaskSchedule(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_scim.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_task_schedule.default", "crontab", "6 */4 * * 2"),
					resource.TestCheckResourceAttr("authentik_task_schedule.default", "paused", "false"),
				),
			},
		},
	})
}

func testAccResourceTaskSchedule(name string) string {
	return fmt.Sprintf(`
resource "authentik_provider_scim" "name" {
  name  = "%[1]s"
  url   = "http://localhost"
  token = "foo"
}

resource "authentik_task_schedule" "default" {
  app_model = "authentik_providers_scim.scimprovider"
  model_id  = authentik_provider_scim.name.id
  actor_name = "authentik.providers.scim.tasks.scim_sync"

  crontab = "6 */4 * * 2"
}
`, name)
}
