package sdkprovider

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fwprovider "goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestProvider(t *testing.T) {
	p := Provider("testing", false)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %[1]s", err)
	}
}

// schemaCmpOptions mirrors tf6muxserver's own (unexported) schema comparison options:
// SchemaAttribute/SchemaNestedBlock slices may differ in order without being a real
// mismatch, and MinItems/MaxItems on nested blocks are excluded from mux's own equality
// check. See github.com/hashicorp/terraform-plugin-mux/tf6muxserver/schema_equality.go.
var schemaCmpOptions = []cmp.Option{
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaAttribute) bool {
		return i.Name < j.Name
	}),
	cmpopts.SortSlices(func(i, j *tfprotov6.SchemaNestedBlock) bool {
		return i.TypeName < j.TypeName
	}),
	cmpopts.IgnoreFields(tfprotov6.SchemaNestedBlock{}, "MinItems", "MaxItems"),
}

// TestProviderSchemaMatchesSDKv2 guards the mux-compatibility work in 0a: if the SDKv2
// provider block schema and the framework provider block schema ever diverge,
// tf6muxserver.NewMuxServer's own comparison would reject the pairing at startup with
// "Invalid Provider Server Combination" - this test catches that at `go test` time
// instead of `terraform plan` time.
func TestProviderSchemaMatchesSDKv2(t *testing.T) {
	ctx := context.Background()

	sdkServer, err := tf5to6server.UpgradeServer(ctx, Provider("testing", false).GRPCProvider)
	require.NoError(t, err)
	sdkSchema, err := sdkServer.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, sdkSchema.Diagnostics)

	fwServer := providerserver.NewProtocol6(fwprovider.New("testing", false))()
	fwSchema, err := fwServer.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, fwSchema.Diagnostics)

	if diff := cmp.Diff(sdkSchema.Provider, fwSchema.Provider, schemaCmpOptions...); diff != "" {
		t.Errorf("provider block schema mismatch between pkg/sdkprovider and pkg/provider (-sdk +framework):\n%s", diff)
	}
}

func TestProviderConfigure_PathBasedURL(t *testing.T) {
	testCases := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "Root path with trailing slash",
			inputURL:    "https://api.example.com/",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Root path without trailing slash",
			inputURL:    "https://api.example.com",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Single segment path with trailing slash",
			inputURL:    "https://api.example.com/sso/",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Single segment path without trailing slash",
			inputURL:    "https://api.example.com/sso",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Multi-segment path with trailing slash",
			inputURL:    "https://api.example.com/auth/sso/",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "Multi-segment path without trailing slash",
			inputURL:    "https://api.example.com/auth/sso",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "HTTP scheme with path",
			inputURL:    "http://localhost:9000/sso/",
			expectedURL: "http://localhost:9000/sso/api/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Provider("testing", true)

			_ac, diag := p.ConfigureContextFunc(t.Context(), schema.TestResourceDataRaw(t, p.Schema, map[string]any{
				"url":      tc.inputURL,
				"token":    "test-token",
				"insecure": false,
			}))
			assert.Nil(t, diag)
			ac := _ac.(*APIClient)
			assert.Equal(t, tc.expectedURL, ac.client.GetConfig().Servers[0].URL, "Server URL should be constructed correctly")
		})
	}
}
