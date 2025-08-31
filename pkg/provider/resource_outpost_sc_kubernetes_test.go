package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceServiceConnectionKubernetes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceConnectionKubernetes(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_service_connection_kubernetes.name", "name", rName),
				),
			},
			{
				Config: testAccResourceServiceConnectionKubernetes(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_service_connection_kubernetes.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceServiceConnectionKubernetes(name string) string {
	return fmt.Sprintf(`
resource "authentik_service_connection_kubernetes" "name" {
  name = "%[1]s"
  local = true
}
`, name)
}
