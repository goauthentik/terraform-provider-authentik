package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceServiceConnectionDocker(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
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
