package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegistrySchemasValidate walks every resource and data source New() registers with
// the framework provider and validates its schema via the framework's own
// ValidateImplementation, table-driven over the registry rather than one test per
// resource. This is what catches H2 (Default set without Computed, and friends) at
// `go test` time instead of the "GetProviderSchema" RPC failing at `terraform plan` time.
func TestRegistrySchemasValidate(t *testing.T) {
	ctx := context.Background()
	p := New("testing", true)

	t.Run("resources", func(t *testing.T) {
		for _, newResource := range p.Resources(ctx) {
			r := newResource()

			var metaResp fwresource.MetadataResponse
			r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "authentik"}, &metaResp)

			t.Run(metaResp.TypeName, func(t *testing.T) {
				var schemaResp fwresource.SchemaResponse
				r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
				require.Empty(t, schemaResp.Diagnostics)

				diags := schemaResp.Schema.ValidateImplementation(ctx)
				assert.Empty(t, diags)
			})
		}
	})

	t.Run("data_sources", func(t *testing.T) {
		for _, newDataSource := range p.DataSources(ctx) {
			d := newDataSource()

			var metaResp datasource.MetadataResponse
			d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "authentik"}, &metaResp)

			t.Run(metaResp.TypeName, func(t *testing.T) {
				var schemaResp datasource.SchemaResponse
				d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
				require.Empty(t, schemaResp.Diagnostics)

				diags := schemaResp.Schema.ValidateImplementation(ctx)
				assert.Empty(t, diags)
			})
		}
	})
}
