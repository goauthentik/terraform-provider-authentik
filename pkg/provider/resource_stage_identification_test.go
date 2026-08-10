package provider_test

import (
	"fmt"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceStageIdentification(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStageIdentification(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_identification.name", "name", rName+"-ident"),
				),
			},
			{
				Config: testAccResourceStageIdentification(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_identification.name", "name", rName+"test-ident"),
				),
			},
		},
	})
}

func testAccResourceStageIdentification(name string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_oauth" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id

  provider_type = "discord"
  consumer_key = "foo"
  consumer_secret = "bar"
}

resource "authentik_stage_password" "name" {
  name              = "%[1]s-pass"
  backends = ["authentik.core.auth.InbuiltBackend"]
}

resource "authentik_stage_identification" "name" {
  name              = "%[1]s-ident"
  user_fields = ["username"]
  sources = [authentik_source_oauth.name.uuid]
  password_stage = authentik_stage_password.name.id
}
`, name)
}
