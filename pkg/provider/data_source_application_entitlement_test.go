package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceApplicationEntitlement(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplicationEntitlementByIdConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_application_entitlement.test", "name", rName),
				),
			},
			{
				Config: testAccDataSourceApplicationEntitlementByNameConfig(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_application_entitlement.test", "name", rName),
				),
			},
		},
	})
}

func TestAccDataSourceApplicationEntitlement_Errors(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceApplicationEntitlementNoConfigSimple,
				ExpectError: regexp.MustCompile(`Neither id nor app/name pair were provided`),
			},
			{
				Config:      testAccDataSourceApplicationEntitlementBrokenConfig1Simple,
				ExpectError: regexp.MustCompile(`Conflicting configuration arguments`),
			},
			{
				Config:      testAccDataSourceApplicationEntitlementBrokenConfig2Simple,
				ExpectError: regexp.MustCompile(`Neither id nor app/name pair were provided`),
			},
			{
				Config:      testAccDataSourceApplicationEntitlementNotFound(appName, rName),
				ExpectError: regexp.MustCompile(`No matching application entitlements found`),
			},
		},
	})
}

func testAccDataSourceApplicationEntitlementByIdConfig(appName string, entitlementName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_application_entitlement" "test" {
  application = authentik_application.test.uuid
  name        = "%[2]s"
}

data "authentik_application_entitlement" "test" {
  id = authentik_application_entitlement.test.id
}
`, appName, entitlementName)
}

func testAccDataSourceApplicationEntitlementByNameConfig(appName string, entitlementName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

resource "authentik_application_entitlement" "test" {
  application = authentik_application.test.uuid
  name        = "%[2]s"
}

data "authentik_application_entitlement" "test" {
  app  = authentik_application_entitlement.test.application
  name = authentik_application_entitlement.test.name
}
`, appName, entitlementName)
}

const testAccDataSourceApplicationEntitlementNoConfigSimple = `
data "authentik_application_entitlement" "missing-config" {
}
`

const testAccDataSourceApplicationEntitlementBrokenConfig1Simple = `
data "authentik_application_entitlement" "broken-config-1" {
  id  = "12345678-1234-1234-1234-123456789012"
  app = "12345678-1234-1234-1234-123456789012"
}
`

const testAccDataSourceApplicationEntitlementBrokenConfig2Simple = `
data "authentik_application_entitlement" "broken-config-2" {
  app = "12345678-1234-1234-1234-123456789012"
}
`

func testAccDataSourceApplicationEntitlementNotFound(appName string, entitlementName string) string {
	return fmt.Sprintf(`
resource "authentik_application" "test" {
  name = "%[1]s"
  slug = "%[1]s"
}

data "authentik_application_entitlement" "test" {
  app  = authentik_application.test.uuid
  name = "%[2]s"
}
`, appName, entitlementName)
}
