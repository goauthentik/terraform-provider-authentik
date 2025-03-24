package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicyExpression(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
