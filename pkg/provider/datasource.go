package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// dataSourceBase is embedded by every framework data source. It is resourceBase minus
// ImportState, which data sources have no equivalent of.
type dataSourceBase struct {
	client *api.APIClient
}

// Configure implements datasource.DataSourceWithConfigure. See resourceBase.Configure
// for why req.ProviderData == nil must be a no-op rather than an error.
func (d *dataSourceBase) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// span starts a Sentry span for a data source Read. Use as
// `defer d.span(ctx, "read")()` at the top of Read.
func (d *dataSourceBase) span(ctx context.Context, operation string) func() {
	return helpers.Span(ctx, "terraform.datasource", "terraform.datasource."+operation, "Datasource "+operation)
}
