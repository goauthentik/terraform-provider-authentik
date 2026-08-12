package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyPassword(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyPassword(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_password.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyPassword(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_password.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyPassword(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_password" "name" {
  name              = "%[1]s"
  error_message = "foo"
}
`, name)
}
