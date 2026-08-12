package provider_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	pkgacctest "goauthentik.io/terraform-provider-authentik/pkg/acctest"
)

func TestAccResourceSourcePlex(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { pkgacctest.PreCheck(t) },
		ProtoV6ProviderFactories: pkgacctest.ProviderFactories,
		CheckDestroy:             testAccCheckSourcePlexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSourcePlex(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_plex.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_source_plex.name", "id", rName),
				),
			},
			// This step is the H3 exercise for the whole sources batch: it changes the slug,
			// so the resource's id changes with it. The plan's id is the unknown placeholder
			// at that point, which is why Update has to take the prior id from req.State -
			// using req.Plan would send the placeholder into the URL. The id check below
			// asserts state actually tracks the new slug.
			{
				Config: testAccResourceSourcePlex(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_source_plex.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_source_plex.name", "id", rName+"test"),
					resource.TestCheckResourceAttr("authentik_source_plex.name", "slug", rName+"test"),
				),
			},
			// plex is the right source to carry this: unlike the other four secret-bearing
			// sources, the API does return plex_token, so nothing needs
			// ImportStateVerifyIgnore and the round trip is exact - including uuid, which is
			// Optional+Computed and pinned with UseStateForUnknown.
			{
				ResourceName:      "authentik_source_plex.name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSourcePlexDestroy(s *terraform.State) error {
	c := pkgacctest.APIClientFromEnv()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "authentik_source_plex" {
			continue
		}
		_, hr, err := c.SourcesAPI.SourcesPlexRetrieve(context.Background(), rs.Primary.ID).Execute()
		if err == nil {
			return fmt.Errorf("plex source %s still exists", rs.Primary.ID)
		}
		if hr == nil || hr.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

func testAccResourceSourcePlex(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_source_plex" "name" {
  name      = "%[1]s"
  slug      = "%[1]s"
  authentication_flow = data.authentik_flow.default-authorization-flow.id
  enrollment_flow = data.authentik_flow.default-authorization-flow.id
  client_id = "foo-bar-baz"
  plex_token = "foo"
}
`, name, appName)
}
