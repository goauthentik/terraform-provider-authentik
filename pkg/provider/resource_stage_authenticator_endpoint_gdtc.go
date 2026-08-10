package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &stageAuthenticatorEndpointGDTCResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorEndpointGDTCResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorEndpointGDTCResource{}
)

func newStageAuthenticatorEndpointGDTCResource() resource.Resource {
	return &stageAuthenticatorEndpointGDTCResource{}
}

type stageAuthenticatorEndpointGDTCResource struct {
	resourceBase
}

type stageAuthenticatorEndpointGDTCModel struct {
	ID            types.String         `tfsdk:"id"`
	Name          types.String         `tfsdk:"name"`
	FriendlyName  types.String         `tfsdk:"friendly_name"`
	ConfigureFlow types.String         `tfsdk:"configure_flow"`
	Credentials   jsontypes.Normalized `tfsdk:"credentials"`
}

func (r *stageAuthenticatorEndpointGDTCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_endpoint_gdtc"
}

func (r *stageAuthenticatorEndpointGDTCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Flows & Stages --- ",
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
			"friendly_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             friendlyNameDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(friendlyNameDefault.Value())),
			},
			"configure_flow": schema.StringAttribute{
				Optional: true,
			},
			// No Description in the SDKv2 schema for this one attribute - not even
			// helpers.JSONDescription, unlike every other JSON-typed attribute.
			"credentials": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Required:   true,
				Sensitive:  true,
			},
		},
	}
}

func (r *stageAuthenticatorEndpointGDTCResource) toRequest(data *stageAuthenticatorEndpointGDTCModel) (*api.AuthenticatorEndpointGDTCStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var credentials map[string]any
	diags.Append(data.Credentials.Unmarshal(&credentials)...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.AuthenticatorEndpointGDTCStageRequest{
		Name:          data.Name.ValueString(),
		ConfigureFlow: *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
		FriendlyName:  new(data.FriendlyName.ValueString()),
		Credentials:   credentials,
	}, diags
}

// fromAPI deliberately does not set FriendlyName/ConfigureFlow (H4): the SDKv2 version
// never wrote them back either, so per the migration-plan.md H4 fix, data being seeded
// from the plan/prior state before this runs means they retain their last-known value
// for free, with no special-casing needed here.
func (r *stageAuthenticatorEndpointGDTCResource) fromAPI(data *stageAuthenticatorEndpointGDTCModel, res *api.AuthenticatorEndpointGDTCStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	credBytes, err := json.Marshal(res.Credentials)
	if err != nil {
		diags.AddError("Failed to encode credentials", err.Error())
		return diags
	}
	data.Credentials = jsontypes.NewNormalizedValue(string(credBytes))

	return diags
}

func (r *stageAuthenticatorEndpointGDTCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorEndpointGDTCModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEndpointGdtcCreate(ctx).AuthenticatorEndpointGDTCStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorEndpointGDTCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorEndpointGDTCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEndpointGdtcRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorEndpointGDTCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorEndpointGDTCModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEndpointGdtcUpdate(ctx, data.ID.ValueString()).AuthenticatorEndpointGDTCStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorEndpointGDTCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorEndpointGDTCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorEndpointGdtcDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
