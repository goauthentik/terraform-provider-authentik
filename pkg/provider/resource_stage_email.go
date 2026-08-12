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
	_ resource.Resource                = &stageEmailResource{}
	_ resource.ResourceWithConfigure   = &stageEmailResource{}
	_ resource.ResourceWithImportState = &stageEmailResource{}
)

func newStageEmailResource() resource.Resource {
	return &stageEmailResource{}
}

type stageEmailResource struct {
	resourceBase
}

type stageEmailModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	UseGlobalSettings     types.Bool   `tfsdk:"use_global_settings"`
	Host                  types.String `tfsdk:"host"`
	Port                  types.Int32  `tfsdk:"port"`
	Username              types.String `tfsdk:"username"`
	Password              types.String `tfsdk:"password"`
	UseTLS                types.Bool   `tfsdk:"use_tls"`
	UseSSL                types.Bool   `tfsdk:"use_ssl"`
	Timeout               types.Int32  `tfsdk:"timeout"`
	FromAddress           types.String `tfsdk:"from_address"`
	TokenExpiry           types.String `tfsdk:"token_expiry"`
	Subject               types.String `tfsdk:"subject"`
	Template              types.String `tfsdk:"template"`
	ActivateUserOnSuccess types.Bool   `tfsdk:"activate_user_on_success"`
	RecoveryMaxAttempts   types.Int32  `tfsdk:"recovery_max_attempts"`
	RecoveryCacheTimeout  types.String `tfsdk:"recovery_cache_timeout"`
}

func (r *stageEmailResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_email"
}

func (r *stageEmailResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	useGlobalSettingsDefault := helpers.BoolDefault(true)
	hostDefault := helpers.StringDefault("localhost")
	portDefault := helpers.Int32Default(25)
	timeoutDefault := helpers.Int32Default(30)
	fromAddressDefault := helpers.StringDefault("system@authentik.local")
	tokenExpiryDefault := helpers.StringDefault("minutes=30")
	subjectDefault := helpers.StringDefault("authentik")
	templateDefault := helpers.StringDefault("email/password_reset.html")
	activateUserOnSuccessDefault := helpers.BoolDefault(false)
	recoveryMaxAttemptsDefault := helpers.Int32Default(5)
	recoveryCacheTimeoutDefault := helpers.StringDefault("minutes=5")

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
			// use_tls and use_ssl are the only two attributes here with no Default, so
			// they are also the only two that keep the prior-aware read helpers.
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
			"activate_user_on_success": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             activateUserOnSuccessDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(activateUserOnSuccessDefault.Value())),
			},
			"recovery_max_attempts": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             recoveryMaxAttemptsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(recoveryMaxAttemptsDefault.Value())),
			},
			"recovery_cache_timeout": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             recoveryCacheTimeoutDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(recoveryCacheTimeoutDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
		},
	}
}

func (r *stageEmailResource) toRequest(data *stageEmailModel) *api.EmailStageRequest {
	return &api.EmailStageRequest{
		Name:              data.Name.ValueString(),
		UseGlobalSettings: new(data.UseGlobalSettings.ValueBool()),
		// use_tls/use_ssl have no Default, but SDKv2 sent new(d.Get(...).(bool))
		// unconditionally, i.e. false when unset. ValueBool() on null is also false, so
		// the wire format is unchanged.
		UseTls:      new(data.UseTLS.ValueBool()),
		UseSsl:      new(data.UseSSL.ValueBool()),
		Host:        helpers.StringPtr(data.Host),
		Username:    helpers.StringPtr(data.Username),
		Password:    helpers.StringPtr(data.Password),
		FromAddress: helpers.StringPtr(data.FromAddress),
		TokenExpiry: helpers.StringPtr(data.TokenExpiry),
		Subject:     helpers.StringPtr(data.Subject),
		Template:    helpers.StringPtr(data.Template),
		Port:        helpers.Int32Ptr(data.Port),
		Timeout:     helpers.Int32Ptr(data.Timeout),
		// SDKv2 built these two with GetP/GetIntP, which return nil for false and 0
		// respectively, so activate_user_on_success could never be turned back off once
		// enabled (the API keeps absent fields unchanged on PUT). Both have Defaults, so
		// they are never null here and the concrete value is always sent.
		ActivateUserOnSuccess: new(data.ActivateUserOnSuccess.ValueBool()),
		RecoveryMaxAttempts:   helpers.Int32Ptr(data.RecoveryMaxAttempts),
		RecoveryCacheTimeout:  helpers.StringPtr(data.RecoveryCacheTimeout),
	}
}

// fromAPI deliberately never assigns Password: this is one of the 18 H4 resources, and
// SDKv2's Read never called SetWrapper for it either. Callers seed data from req.Plan or
// req.State first, so the omission preserves the configured secret rather than nulling it
// out - see stageCaptchaResource.fromAPI for the full explanation.
func (r *stageEmailResource) fromAPI(data *stageEmailModel, res *api.EmailStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	// No Default on these three, so they keep the prior-aware helpers.
	data.Username = helpers.StringOrNull(data.Username, res.GetUsername())
	data.UseTLS = helpers.BoolOrNull(data.UseTLS, res.GetUseTls())
	data.UseSSL = helpers.BoolOrNull(data.UseSSL, res.GetUseSsl())
	// The rest all have Defaults and take the API value verbatim.
	data.UseGlobalSettings = types.BoolValue(res.GetUseGlobalSettings())
	data.Host = types.StringValue(res.GetHost())
	data.Port = types.Int32Value(res.GetPort())
	data.Timeout = types.Int32Value(res.GetTimeout())
	data.FromAddress = types.StringValue(res.GetFromAddress())
	data.TokenExpiry = types.StringValue(res.GetTokenExpiry())
	data.Subject = types.StringValue(res.GetSubject())
	data.Template = types.StringValue(res.GetTemplate())
	data.ActivateUserOnSuccess = types.BoolValue(res.GetActivateUserOnSuccess())
	data.RecoveryMaxAttempts = types.Int32Value(res.GetRecoveryMaxAttempts())
	data.RecoveryCacheTimeout = types.StringValue(res.GetRecoveryCacheTimeout())
}

func (r *stageEmailResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageEmailModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEmailCreate(ctx).EmailStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageEmailResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries password across a Read - see fromAPI.
	var data stageEmailModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEmailRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageEmailResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageEmailModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEmailUpdate(ctx, data.ID.ValueString()).EmailStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageEmailResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageEmailModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesEmailDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
