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
	_ resource.Resource                = &stageAuthenticatorEmailResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorEmailResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorEmailResource{}
)

func newStageAuthenticatorEmailResource() resource.Resource {
	return &stageAuthenticatorEmailResource{}
}

type stageAuthenticatorEmailResource struct {
	resourceBase
}

type stageAuthenticatorEmailModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	FriendlyName      types.String `tfsdk:"friendly_name"`
	ConfigureFlow     types.String `tfsdk:"configure_flow"`
	UseGlobalSettings types.Bool   `tfsdk:"use_global_settings"`
	Host              types.String `tfsdk:"host"`
	Port              types.Int32  `tfsdk:"port"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	UseTLS            types.Bool   `tfsdk:"use_tls"`
	UseSSL            types.Bool   `tfsdk:"use_ssl"`
	Timeout           types.Int32  `tfsdk:"timeout"`
	FromAddress       types.String `tfsdk:"from_address"`
	TokenExpiry       types.String `tfsdk:"token_expiry"`
	Subject           types.String `tfsdk:"subject"`
	Template          types.String `tfsdk:"template"`
}

func (r *stageAuthenticatorEmailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_email"
}

func (r *stageAuthenticatorEmailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")
	useGlobalSettingsDefault := helpers.BoolDefault(true)
	hostDefault := helpers.StringDefault("localhost")
	portDefault := helpers.Int32Default(25)
	timeoutDefault := helpers.Int32Default(30)
	fromAddressDefault := helpers.StringDefault("system@authentik.local")
	tokenExpiryDefault := helpers.StringDefault("minutes=30")
	subjectDefault := helpers.StringDefault("authentik")
	templateDefault := helpers.StringDefault("email/password_reset.html")

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
			"use_global_settings": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             useGlobalSettingsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(useGlobalSettingsDefault.Value())),
			},
			"host": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             hostDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(hostDefault.Value())),
			},
			"port": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             portDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(portDefault.Value())),
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			// As in authentik_stage_email, these two are the only attributes without a
			// Default and so the only ones keeping the prior-aware read helpers.
			"use_tls": schema.BoolAttribute{
				Optional: true,
			},
			"use_ssl": schema.BoolAttribute{
				Optional: true,
			},
			"timeout": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             timeoutDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(timeoutDefault.Value())),
			},
			"from_address": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             fromAddressDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(fromAddressDefault.Value())),
			},
			"token_expiry": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             tokenExpiryDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(tokenExpiryDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"subject": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             subjectDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(subjectDefault.Value())),
			},
			"template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             templateDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(templateDefault.Value())),
			},
		},
	}
}

func (r *stageAuthenticatorEmailResource) toRequest(data *stageAuthenticatorEmailModel) *api.AuthenticatorEmailStageRequest {
	return &api.AuthenticatorEmailStageRequest{
		Name:              data.Name.ValueString(),
		UseGlobalSettings: new(data.UseGlobalSettings.ValueBool()),
		// No Default on these two, but SDKv2 sent them unconditionally as false when
		// unset, and ValueBool() on null is also false.
		UseTls:        new(data.UseTLS.ValueBool()),
		UseSsl:        new(data.UseSSL.ValueBool()),
		FriendlyName:  helpers.StringPtr(data.FriendlyName),
		ConfigureFlow: *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
		Host:          helpers.StringPtr(data.Host),
		Port:          helpers.Int32Ptr(data.Port),
		Username:      helpers.StringPtr(data.Username),
		Password:      helpers.StringPtr(data.Password),
		Timeout:       helpers.Int32Ptr(data.Timeout),
		FromAddress:   helpers.StringPtr(data.FromAddress),
		TokenExpiry:   helpers.StringPtr(data.TokenExpiry),
		Subject:       helpers.StringPtr(data.Subject),
		Template:      helpers.StringPtr(data.Template),
	}
}

// fromAPI deliberately never assigns Password: this is one of the 18 H4 resources, and
// SDKv2's Read never called SetWrapper for it either. Callers seed data from req.Plan or
// req.State first, so the omission preserves the configured secret rather than nulling it
// out - see stageCaptchaResource.fromAPI for the full explanation.
func (r *stageAuthenticatorEmailResource) fromAPI(data *stageAuthenticatorEmailModel, res *api.AuthenticatorEmailStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
	// No Defaults on these three, so they keep the prior-aware helpers.
	data.Username = helpers.StringOrNull(data.Username, res.GetUsername())
	data.UseTLS = helpers.BoolOrNull(data.UseTLS, res.GetUseTls())
	data.UseSSL = helpers.BoolOrNull(data.UseSSL, res.GetUseSsl())
	// The rest all have Defaults and take the API value verbatim.
	data.FriendlyName = types.StringValue(res.GetFriendlyName())
	data.UseGlobalSettings = types.BoolValue(res.GetUseGlobalSettings())
	data.Host = types.StringValue(res.GetHost())
	data.Port = types.Int32Value(res.GetPort())
	data.Timeout = types.Int32Value(res.GetTimeout())
	data.FromAddress = types.StringValue(res.GetFromAddress())
	data.TokenExpiry = types.StringValue(res.GetTokenExpiry())
	data.Subject = types.StringValue(res.GetSubject())
	data.Template = types.StringValue(res.GetTemplate())
}

func (r *stageAuthenticatorEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorEmailModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEmailCreate(ctx).AuthenticatorEmailStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries password across a Read - see fromAPI.
	var data stageAuthenticatorEmailModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEmailRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorEmailModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorEmailUpdate(ctx, data.ID.ValueString()).AuthenticatorEmailStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorEmailModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorEmailDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
