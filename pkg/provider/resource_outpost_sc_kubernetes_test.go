package provider

import (
	"fmt"
	"regexp"
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
			{
				Config:      testAccResourceServiceConnectionKubernetesKubeconfig(rName),
				ExpectError: regexp.MustCompile(`invalid\scharacter`),
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

func testAccResourceServiceConnectionKubernetesKubeconfig(name string) string {
	return fmt.Sprintf(`
resource "authentik_service_connection_kubernetes" "name" {
  name = "%[1]s"
  kubeconfig = "foo"
}
`, name)
}
