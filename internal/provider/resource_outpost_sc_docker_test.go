package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceServiceConnectionDocker(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceConnectionDocker(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_service_connection_docker.name", "name", rName),
				),
			},
			{
				Config: testAccResourceServiceConnectionDocker(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_service_connection_docker.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceServiceConnectionDocker(name string) string {
	return fmt.Sprintf(`
resource "authentik_service_connection_docker" "name" {
  name = "%[1]s"
  local = true
}
`, name)
}
