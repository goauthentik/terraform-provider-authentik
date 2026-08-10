package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyExpression(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyExpression(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_expression.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyExpression(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_expression.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyExpression(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_expression" "name" {
  name              = "%[1]s"
  expression = "return True"
}
resource "authentik_policy_expression" "name2" {
  name              = "%[1]s-EOT"
  expression = <<EOT
return True
EOT
}
`, name)
}
