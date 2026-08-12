package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &stageAuthenticatorDuoResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorDuoResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorDuoResource{}
)

func newStageAuthenticatorDuoResource() resource.Resource {
	return &stageAuthenticatorDuoResource{}
}

type stageAuthenticatorDuoResource struct {
	resourceBase
}

type stageAuthenticatorDuoModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	FriendlyName        types.String `tfsdk:"friendly_name"`
	ConfigureFlow       types.String `tfsdk:"configure_flow"`
	ClientID            types.String `tfsdk:"client_id"`
	ClientSecret        types.String `tfsdk:"client_secret"`
	AdminIntegrationKey types.String `tfsdk:"admin_integration_key"`
	AdminSecretKey      types.String `tfsdk:"admin_secret_key"`
	APIHostname         types.String `tfsdk:"api_hostname"`
}

func (r *stageAuthenticatorDuoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_duo"
}

func (r *stageAuthenticatorDuoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"client_id": schema.StringAttribute{
				Required: true,
			},
			"client_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"admin_integration_key": schema.StringAttribute{
				Optional: true,
			},
			"admin_secret_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"api_hostname": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *stageAuthenticatorDuoResource) toRequest(data *stageAuthenticatorDuoModel) *api.AuthenticatorDuoStageRequest {
	return &api.AuthenticatorDuoStageRequest{
		Name:                data.Name.ValueString(),
		ClientId:            data.ClientID.ValueString(),
		ClientSecret:        data.ClientSecret.ValueString(),
		ApiHostname:         data.APIHostname.ValueString(),
		FriendlyName:        helpers.StringPtr(data.FriendlyName),
		AdminIntegrationKey: helpers.StringPtr(data.AdminIntegrationKey),
		AdminSecretKey:      helpers.StringPtr(data.AdminSecretKey),
		ConfigureFlow:       *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
	}
}

// fromAPI deliberately never assigns ClientSecret or AdminSecretKey: this is one of the 18
// H4 resources, and the API returns neither. Callers seed data from req.Plan or req.State
// first, so the omissions preserve the configured secrets rather than nulling them out -
// see stageCaptchaResource.fromAPI for the full explanation. client_secret is Required, so
// nulling it would fail the apply outright rather than merely losing a value.
func (r *stageAuthenticatorDuoResource) fromAPI(data *stageAuthenticatorDuoModel, res *api.AuthenticatorDuoStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ClientID = types.StringValue(res.ClientId)
	data.APIHostname = types.StringValue(res.ApiHostname)
	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
	// No Default, so prior-aware.
	data.AdminIntegrationKey = helpers.StringOrNull(data.AdminIntegrationKey, res.GetAdminIntegrationKey())
	// Has a Default, so verbatim.
	data.FriendlyName = types.StringValue(res.GetFriendlyName())
}

func (r *stageAuthenticatorDuoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorDuoModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorDuoCreate(ctx).AuthenticatorDuoStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorDuoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries the two secrets across a Read - see fromAPI.
	var data stageAuthenticatorDuoModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorDuoRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorDuoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorDuoModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorDuoUpdate(ctx, data.ID.ValueString()).AuthenticatorDuoStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorDuoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorDuoModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorDuoDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
