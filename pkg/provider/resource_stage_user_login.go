package provider

import (
	"context"

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
	_ resource.Resource                = &stageUserLoginResource{}
	_ resource.ResourceWithConfigure   = &stageUserLoginResource{}
	_ resource.ResourceWithImportState = &stageUserLoginResource{}
)

func newStageUserLoginResource() resource.Resource {
	return &stageUserLoginResource{}
}

type stageUserLoginResource struct {
	resourceBase
}

type stageUserLoginModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	SessionDuration        types.String `tfsdk:"session_duration"`
	RememberMeOffset       types.String `tfsdk:"remember_me_offset"`
	TerminateOtherSessions types.Bool   `tfsdk:"terminate_other_sessions"`
	NetworkBinding         types.String `tfsdk:"network_binding"`
	GeoipBinding           types.String `tfsdk:"geoip_binding"`
	RememberDevice         types.String `tfsdk:"remember_device"`
}

func (r *stageUserLoginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_user_login"
}

func (r *stageUserLoginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	sessionDurationDefault := helpers.StringDefault("seconds=0")
	rememberMeOffsetDefault := helpers.StringDefault("seconds=0")
	terminateOtherSessionsDefault := helpers.BoolDefault(false)
	networkBindingDefault := helpers.StringDefault(string(api.NETWORKBINDINGENUM_NO_BINDING))
	geoipBindingDefault := helpers.StringDefault(string(api.GEOIPBINDINGENUM_NO_BINDING))
	rememberDeviceDefault := helpers.StringDefault("days=30")

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
			"session_duration": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             sessionDurationDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(sessionDurationDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"remember_me_offset": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             rememberMeOffsetDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(rememberMeOffsetDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"terminate_other_sessions": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             terminateOtherSessionsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(terminateOtherSessionsDefault.Value())),
			},
			"network_binding": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             networkBindingDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedNetworkBindingEnumEnumValues), helpers.WithDefault(networkBindingDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedNetworkBindingEnumEnumValues),
				},
			},
			"geoip_binding": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             geoipBindingDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedGeoipBindingEnumEnumValues), helpers.WithDefault(geoipBindingDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedGeoipBindingEnumEnumValues),
				},
			},
			"remember_device": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             rememberDeviceDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(rememberDeviceDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
		},
	}
}

func (r *stageUserLoginResource) toRequest(data *stageUserLoginModel) *api.UserLoginStageRequest {
	return &api.UserLoginStageRequest{
		Name:                   data.Name.ValueString(),
		SessionDuration:        new(data.SessionDuration.ValueString()),
		TerminateOtherSessions: new(data.TerminateOtherSessions.ValueBool()),
		RememberMeOffset:       new(data.RememberMeOffset.ValueString()),
		NetworkBinding:         api.NetworkBindingEnum(data.NetworkBinding.ValueString()).Ptr(),
		GeoipBinding:           api.GeoipBindingEnum(data.GeoipBinding.ValueString()).Ptr(),
		RememberDevice:         new(data.RememberDevice.ValueString()),
	}
}

// fromAPI deliberately never assigns NetworkBinding or GeoipBinding. This is one of the 18
// H4 resources, and it is the subtler kind: unlike the write-only secrets in part 6, the
// API *does* return both of these fields - SDKv2's Read simply never called SetWrapper for
// them, so they were never refreshed from the server and state always kept the configured
// value. Callers seed the model from req.Plan or req.State before calling this, so the
// omission reproduces that exactly.
//
// Note this overrides discovery #5 for these two attributes: they have Defaults, so the
// rule would otherwise say "map verbatim", but H4 says do not map them at all. H4 wins -
// the same call was made for authentik_stage_authenticator_endpoint_gdtc's friendly_name
// in part 3. The practical consequence is that a drift in either field is invisible to
// terraform plan, which is pre-existing behaviour rather than something introduced here.
func (r *stageUserLoginResource) fromAPI(data *stageUserLoginModel, res *api.UserLoginStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	// All four of these have Defaults and are written by SDKv2's Read, so they take the
	// API value verbatim.
	data.SessionDuration = types.StringValue(res.GetSessionDuration())
	data.RememberMeOffset = types.StringValue(res.GetRememberMeOffset())
	data.TerminateOtherSessions = types.BoolValue(res.GetTerminateOtherSessions())
	data.RememberDevice = types.StringValue(res.GetRememberDevice())
}

func (r *stageUserLoginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageUserLoginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLoginCreate(ctx).UserLoginStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserLoginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries network_binding/geoip_binding across a
	// Read - see fromAPI.
	var data stageUserLoginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLoginRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageUserLoginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageUserLoginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLoginUpdate(ctx, data.ID.ValueString()).UserLoginStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserLoginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageUserLoginModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesUserLoginDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
