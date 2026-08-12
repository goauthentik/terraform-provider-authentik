package sdkprovider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceToken(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	expires := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceToken(rName, expires),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_token.token", "identifier", rName),
					resource.TestCheckResourceAttrSet("authentik_token.token", "key"),
				),
			},
		},
	})
}

func testAccResourceToken(name string, time string) string {
	return fmt.Sprintf(`
resource "authentik_user" "name" {
	username = "%[1]s"
	name = "%[1]s"
}

resource "authentik_token" "token" {
	user = authentik_user.name.id
	identifier = "%[1]s"
	expires = "%[2]s"
	description = "%[1]s"
	retrieve_key = true
}
`, name, time)
}
