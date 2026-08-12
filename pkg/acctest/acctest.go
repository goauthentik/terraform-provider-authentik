// Package acctest holds acceptance-test scaffolding shared between pkg/sdkprovider and
// pkg/provider test files. It exists because a migrated resource's test moves packages
// along with its implementation (see migration-plan.md Phase 2/3), but every acceptance
// test - regardless of which of the two provider packages implements the resource under
// test - needs to exercise the same muxed protocol 6 server main.go actually serves.
package acctest

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	testinginterface "github.com/mitchellh/go-testing-interface"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
	"goauthentik.io/terraform-provider-authentik/pkg/sdkprovider"
)

// MuxServer builds the same muxed protocol 6 server main.go serves, so acceptance tests
// exercise the exact composition users get regardless of which underlying package
// implements the resource or data source under test.
func MuxServer(version string, testing bool) (tfprotov6.ProviderServer, error) {
	upgradedSDKServer, err := tf5to6server.UpgradeServer(context.Background(), sdkprovider.Provider(version, testing).GRPCProvider)
	if err != nil {
		return nil, err
	}

	return tf6muxserver.NewMuxServer(context.Background(),
		func() tfprotov6.ProviderServer { return upgradedSDKServer },
		providerserver.NewProtocol6(provider.New(version, testing)),
	)
}

// ProviderFactories are used to instantiate the muxed provider during acceptance
// testing. The factory function will be invoked for every Terraform CLI command
// executed to create a provider server to which the CLI can reattach.
var ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"authentik": func() (tfprotov6.ProviderServer, error) {
		return MuxServer("test", false)
	},
}

// ProviderTestFactories is ProviderFactories with testing=true: every HTTP request
// fails with a mocked 400 response instead of reaching a real server. Used by tests
// that only need a live-looking provider to configure, not to actually call the API.
var ProviderTestFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"authentik": func() (tfprotov6.ProviderServer, error) {
		return MuxServer("test", true)
	},
}

// PreCheck is the standard resource.TestCase.PreCheck for every acceptance test: it
// fails fast with a clear message if AUTHENTIK_URL/AUTHENTIK_TOKEN aren't set, rather
// than letting the test fail deep inside a Configure error.
func PreCheck(t *testing.T) {
	envIsSet("AUTHENTIK_URL", t)
	envIsSet("AUTHENTIK_TOKEN", t)
}

func envIsSet(k string, t *testing.T) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("%s must be set for acceptance tests", k)
	}
}

// APIClientFromEnv builds a live API client from AUTHENTIK_URL/AUTHENTIK_TOKEN, for use
// in CheckDestroy functions, which only receive a *terraform.State and have no
// *testing.T to plumb through. Goes through the SDKv2 provider's own
// ConfigureContextFunc (rather than helpers.NewAPIClient directly) so it exercises the
// same env-fallback path a real `terraform apply` would.
func APIClientFromEnv() *api.APIClient {
	rt := &testinginterface.RuntimeT{}
	p := sdkprovider.Provider("test", false)
	raw, diags := p.ConfigureContextFunc(context.Background(), schema.TestResourceDataRaw(rt, p.Schema, map[string]any{}))
	if diags.HasError() {
		rt.Fatalf("failed to configure provider: %v", diags)
	}
	return raw.(*sdkprovider.APIClient).Client()
}
