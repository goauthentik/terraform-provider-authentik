package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &providerSSFResource{}
	_ resource.ResourceWithConfigure   = &providerSSFResource{}
	_ resource.ResourceWithImportState = &providerSSFResource{}
)

func newProviderSSFResource() resource.Resource {
	return &providerSSFResource{}
}

type providerSSFResource struct {
	resourceBase
}

type providerSSFModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SigningKey             types.String `tfsdk:"signing_key"`
	JWTFederationProviders types.List   `tfsdk:"jwt_federation_providers"`
	EventRetention         types.String `tfsdk:"event_retention"`
	PushVerifyCertificates types.Bool   `tfsdk:"push_verify_certificates"`
}

func (r *providerSSFResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_ssf"
}

func (r *providerSSFResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	eventRetentionDefault := helpers.StringDefault("days=30")
	pushVerifyCertificatesDefault := helpers.BoolDefault(true)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Applications --- ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			// signing_key, jwt_federation_providers and event_retention are three of the 18
			// H4 attributes: SDKv2's Read wrote only name and push_verify_certificates, so
			// none of the three was ever refreshed from the server. See fromAPI.
			"signing_key": schema.StringAttribute{
				Optional: true,
			},
			"jwt_federation_providers": schema.ListAttribute{
				ElementType:         types.Int32Type,
				Optional:            true,
				MarkdownDescription: "JWTs issued by any of the configured providers can be used to authenticate on behalf of this provider.",
			},
			"event_retention": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             eventRetentionDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(eventRetentionDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"push_verify_certificates": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             pushVerifyCertificatesDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(pushVerifyCertificatesDefault.Value())),
			},
		},
	}
}

func (r *providerSSFResource) toRequest(ctx context.Context, data *providerSSFModel) (*api.SSFProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The Terraform attribute is jwt_federation_providers; the wire field is
	// oidc_auth_providers.
	oidcAuthProviders, d := helpers.SliceOrEmpty[int32](ctx, data.JWTFederationProviders)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.SSFProviderRequest{
		Name:                   data.Name.ValueString(),
		SigningKey:             data.SigningKey.ValueString(),
		EventRetention:         new(data.EventRetention.ValueString()),
		PushVerifyCertificates: new(data.PushVerifyCertificates.ValueBool()),
		OidcAuthProviders:      oidcAuthProviders,
	}, diags
}

// fromAPI deliberately assigns only name and push_verify_certificates. This is one of the 18
// H4 resources and the widest of them: SDKv2's Read never called SetWrapper for
// signing_key, jwt_federation_providers or event_retention, so all three kept whatever was
// already in state and were never refreshed from the server. Callers seed data from req.Plan
// or req.State first, so leaving them unassigned reproduces that exactly - see
// stageCaptchaResource.fromAPI for the general explanation.
//
// Note event_retention has a schema Default, which would normally mean "map verbatim"
// (discovery #5), but H4 wins where the two collide: the attribute is not refreshed at all.
// Same call as for authentik_stage_user_login's bindings.
func (r *providerSSFResource) fromAPI(data *providerSSFModel, res *api.SSFProvider) {
	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.PushVerifyCertificates = types.BoolValue(res.GetPushVerifyCertificates())
}

func (r *providerSSFResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerSSFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersSsfCreate(ctx).SSFProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerSSFResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries the three unrefreshed attributes - see
	// fromAPI.
	var data providerSSFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersSsfRetrieve(ctx, id).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerSSFResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerSSFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersSsfUpdate(ctx, id).SSFProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerSSFResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerSSFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersSsfDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
