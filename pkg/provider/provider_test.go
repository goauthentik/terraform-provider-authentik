package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"authentik": func() (*schema.Provider, error) {
		return Provider("test", false), nil
	},
}

var providerTestFactories = map[string]func() (*schema.Provider, error){
	"authentik": func() (*schema.Provider, error) {
		return Provider("test", true), nil
	},
}

func TestProvider(t *testing.T) {
	p := Provider("testing", false)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %[1]s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	testEnvIsSet("AUTHENTIK_URL", t)
	testEnvIsSet("AUTHENTIK_TOKEN", t)
}

func testEnvIsSet(k string, t *testing.T) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("%[1]s must be set for acceptance tests", k)
	}
}
