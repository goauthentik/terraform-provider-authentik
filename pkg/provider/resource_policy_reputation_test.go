package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePolicyReputation(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyReputation(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_reputation.name", "name", rName),
				),
			},
			{
				Config: testAccResourcePolicyReputation(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_policy_reputation.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourcePolicyReputation(name string) string {
	return fmt.Sprintf(`
resource "authentik_policy_reputation" "name" {
  name              = "%[1]s"
}
`, name)
}
