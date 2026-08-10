package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// resourceBase is embedded by every framework resource. It provides Configure and
// ImportState so each resource only implements Metadata/Schema/Create/Read/
// Update/Delete - the 96 resources being migrated all use plain
// ImportStatePassthroughContext today, so a shared, uniform ImportState is enough.
//
// Sentry tracing deliberately is NOT a wrapper type around resource.Resource the way
// SDKv2's tr() wraps schema.Resource: resource.Resource has nine optional companion
// interfaces (ResourceWithConfigure, ResourceWithImportState,
// ResourceWithConfigValidators, ResourceWithValidateConfig, ResourceWithModifyPlan,
// ResourceWithUpgradeState, ResourceWithMoveState, ResourceWithIdentity,
// ResourceWithUpgradeIdentity), and a wrapper type either implements one or it does not
// - Go has no way to make that conditional on the wrapped value. Embedding resourceBase
// and calling its span method directly is exactly equivalent to today's tracing (tr()
// does not propagate the span's context into the wrapped call either), without that
// hazard.
type resourceBase struct {
	client *api.APIClient
}

// Configure implements resource.ResourceWithConfigure. Per the framework's Configure
// docs, req.ProviderData is nil during the initial "validate" pass before the provider
// itself has been configured, and must be handled as a no-op rather than an error.
func (r *resourceBase) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// ImportState implements resource.ResourceWithImportState via plain ID passthrough,
// matching every one of the 96 SDKv2 resources' ImportStatePassthroughContext today.
func (r *resourceBase) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// span starts a Sentry span for a resource CRUD method. Use as
// `defer r.span(ctx, "create")()` at the top of Create/Read/Update/Delete.
func (r *resourceBase) span(ctx context.Context, operation string) func() {
	return helpers.Span(ctx, "terraform.resource", "terraform.resource."+operation, "Resource "+operation)
}
