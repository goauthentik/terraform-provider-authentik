package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceToken(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	expires := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
