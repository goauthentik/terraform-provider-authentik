package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource                = &stageAuthenticatorValidateResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorValidateResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorValidateResource{}
)

func newStageAuthenticatorValidateResource() resource.Resource {
	return &stageAuthenticatorValidateResource{}
}

type stageAuthenticatorValidateResource struct {
	resourceBase
}

type stageAuthenticatorValidateModel struct {
	ID                         types.String  `tfsdk:"id"`
	Name                       types.String  `tfsdk:"name"`
	NotConfiguredAction        types.String  `tfsdk:"not_configured_action"`
	DeviceClasses              types.List    `tfsdk:"device_classes"`
	ConfigurationStages        types.List    `tfsdk:"configuration_stages"`
	LastAuthThreshold          types.String  `tfsdk:"last_auth_threshold"`
	WebauthnUserVerification   types.String  `tfsdk:"webauthn_user_verification"`
	WebauthnAllowedDeviceTypes types.List    `tfsdk:"webauthn_allowed_device_types"`
	WebauthnHints              types.List    `tfsdk:"webauthn_hints"`
	EmailOtpThrottlingFactor   types.Float64 `tfsdk:"email_otp_throttling_factor"`
	SmsOtpThrottlingFactor     types.Float64 `tfsdk:"sms_otp_throttling_factor"`
	TotpOtpThrottlingFactor    types.Float64 `tfsdk:"totp_otp_throttling_factor"`
	StaticOtpThrottlingFactor  types.Float64 `tfsdk:"static_otp_throttling_factor"`
}

func (r *stageAuthenticatorValidateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_validate"
}

func (r *stageAuthenticatorValidateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	lastAuthThresholdDefault := helpers.StringDefault("seconds=0")
	webauthnUserVerificationDefault := helpers.StringDefault(string(api.USERVERIFICATIONENUM_PREFERRED))
	throttlingDefault := helpers.Float64Default(1)

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
			"not_configured_action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedNotConfiguredActionEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedNotConfiguredActionEnumEnumValues),
				},
			},
			// The four enum lists below carry their validation via
			// listvalidator.ValueStringsAre and deliberately no MarkdownDescription:
			// SDKv2 declared the enum prose on the Elem schema, which CoreConfigSchema
			// drops for primitive lists.
			"device_classes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedDeviceClassesEnumEnumValues)),
				},
			},
			"configuration_stages": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"last_auth_threshold": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             lastAuthThresholdDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(lastAuthThresholdDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"webauthn_user_verification": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             webauthnUserVerificationDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserVerificationEnumEnumValues), helpers.WithDefault(webauthnUserVerificationDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserVerificationEnumEnumValues),
				},
			},
			"webauthn_allowed_device_types": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"webauthn_hints": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedWebAuthnHintEnumEnumValues)),
				},
			},
			"email_otp_throttling_factor": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             throttlingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(throttlingDefault.Value())),
			},
			"sms_otp_throttling_factor": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             throttlingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(throttlingDefault.Value())),
			},
			"totp_otp_throttling_factor": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             throttlingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(throttlingDefault.Value())),
			},
			"static_otp_throttling_factor": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             throttlingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(throttlingDefault.Value())),
			},
		},
	}
}

func (r *stageAuthenticatorValidateResource) toRequest(ctx context.Context, data *stageAuthenticatorValidateModel) (*api.AuthenticatorValidateStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	deviceClasses, d := helpers.SliceOrEmpty[string](ctx, data.DeviceClasses)
	diags.Append(d...)

	configurationStages, d := helpers.SliceOrEmpty[string](ctx, data.ConfigurationStages)
	diags.Append(d...)

	webauthnAllowedDeviceTypes, d := helpers.SliceOrEmpty[string](ctx, data.WebauthnAllowedDeviceTypes)
	diags.Append(d...)

	webauthnHints, d := helpers.SliceOrEmpty[string](ctx, data.WebauthnHints)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.AuthenticatorValidateStageRequest{
		Name: data.Name.ValueString(),
		// #935: not_configured_action is Required and must never be omitted from the
		// request body. It used to be sourced through d.GetOk, which returned a nil
		// pointer for the zero value, and the model's omitempty then dropped the field
		// from the UPDATE PUT entirely, leaving authentik on its own default. Building
		// the pointer unconditionally is the fix - the attribute is Required so it can
		// never be null here, but the point is that nothing in this expression can
		// produce nil. TestStageAuthenticatorValidateToRequest_NotConfiguredAction is
		// the executable spec, ported from the SDKv2 test.
		NotConfiguredAction:        api.NotConfiguredActionEnum(data.NotConfiguredAction.ValueString()).Ptr(),
		DeviceClasses:              helpers.CastSliceString[api.DeviceClassesEnum](deviceClasses),
		ConfigurationStages:        configurationStages,
		WebauthnAllowedDeviceTypes: webauthnAllowedDeviceTypes,
		WebauthnHints:              helpers.CastSliceString[api.WebAuthnHintEnum](webauthnHints),
		LastAuthThreshold:          new(data.LastAuthThreshold.ValueString()),
		WebauthnUserVerification:   api.UserVerificationEnum(data.WebauthnUserVerification.ValueString()).Ptr(),
		EmailOtpThrottlingFactor:   new(data.EmailOtpThrottlingFactor.ValueFloat64()),
		SmsOtpThrottlingFactor:     new(data.SmsOtpThrottlingFactor.ValueFloat64()),
		TotpOtpThrottlingFactor:    new(data.TotpOtpThrottlingFactor.ValueFloat64()),
		StaticOtpThrottlingFactor:  new(data.StaticOtpThrottlingFactor.ValueFloat64()),
	}, diags
}

func (r *stageAuthenticatorValidateResource) fromAPI(ctx context.Context, data *stageAuthenticatorValidateModel, res *api.AuthenticatorValidateStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.NotConfiguredAction = types.StringValue(string(res.GetNotConfiguredAction()))

	// All four lists are Optional-only, so state has to match config exactly after apply.
	// SDKv2 merged configuration_stages and webauthn_allowed_device_types but wrote the
	// API's order straight into state for device_classes and webauthn_hints; all four
	// merge here, which preserves the configured order and is a no-op when the orders
	// already agree.
	deviceClasses := make([]string, len(res.DeviceClasses))
	for i, c := range res.DeviceClasses {
		deviceClasses[i] = string(c)
	}
	mergedDeviceClasses, d := helpers.MergeStringList(ctx, data.DeviceClasses, deviceClasses)
	diags.Append(d...)
	data.DeviceClasses = mergedDeviceClasses

	configurationStages, d := helpers.MergeStringList(ctx, data.ConfigurationStages, res.ConfigurationStages)
	diags.Append(d...)
	data.ConfigurationStages = configurationStages

	webauthnAllowedDeviceTypes, d := helpers.MergeStringList(ctx, data.WebauthnAllowedDeviceTypes, res.WebauthnAllowedDeviceTypes)
	diags.Append(d...)
	data.WebauthnAllowedDeviceTypes = webauthnAllowedDeviceTypes

	webauthnHints := make([]string, len(res.WebauthnHints))
	for i, h := range res.WebauthnHints {
		webauthnHints[i] = string(h)
	}
	mergedWebauthnHints, d := helpers.MergeStringList(ctx, data.WebauthnHints, webauthnHints)
	diags.Append(d...)
	data.WebauthnHints = mergedWebauthnHints

	// Everything below has a Default, so it takes the API value verbatim.
	data.LastAuthThreshold = types.StringValue(res.GetLastAuthThreshold())
	data.WebauthnUserVerification = types.StringValue(string(res.GetWebauthnUserVerification()))
	data.EmailOtpThrottlingFactor = types.Float64Value(res.GetEmailOtpThrottlingFactor())
	data.SmsOtpThrottlingFactor = types.Float64Value(res.GetSmsOtpThrottlingFactor())
	data.TotpOtpThrottlingFactor = types.Float64Value(res.GetTotpOtpThrottlingFactor())
	data.StaticOtpThrottlingFactor = types.Float64Value(res.GetStaticOtpThrottlingFactor())

	return diags
}

func (r *stageAuthenticatorValidateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorValidateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorValidateCreate(ctx).AuthenticatorValidateStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorValidateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorValidateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorValidateRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorValidateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorValidateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorValidateUpdate(ctx, data.ID.ValueString()).AuthenticatorValidateStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorValidateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorValidateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorValidateDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
