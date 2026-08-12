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
	_ resource.Resource                = &stageAuthenticatorWebAuthnResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorWebAuthnResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorWebAuthnResource{}
)

func newStageAuthenticatorWebAuthnResource() resource.Resource {
	return &stageAuthenticatorWebAuthnResource{}
}

type stageAuthenticatorWebAuthnResource struct {
	resourceBase
}

type stageAuthenticatorWebAuthnModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	FriendlyName            types.String `tfsdk:"friendly_name"`
	ConfigureFlow           types.String `tfsdk:"configure_flow"`
	UserVerification        types.String `tfsdk:"user_verification"`
	ResidentKeyRequirement  types.String `tfsdk:"resident_key_requirement"`
	AuthenticatorAttachment types.String `tfsdk:"authenticator_attachment"`
	DeviceTypeRestrictions  types.List   `tfsdk:"device_type_restrictions"`
	MaxAttempts             types.Int32  `tfsdk:"max_attempts"`
	Hints                   types.List   `tfsdk:"hints"`
	PreventDuplicateDevices types.Bool   `tfsdk:"prevent_duplicate_devices"`
}

func (r *stageAuthenticatorWebAuthnResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_webauthn"
}

func (r *stageAuthenticatorWebAuthnResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")
	// Both of these genuinely default to a UserVerificationEnum value in SDKv2, including
	// resident_key_requirement, whose schema reuses the UserVerification enum rather than
	// declaring one of its own. Kept as-is so the docs and the validator stay identical.
	userVerificationDefault := helpers.StringDefault(string(api.USERVERIFICATIONENUM_PREFERRED))
	residentKeyRequirementDefault := helpers.StringDefault(string(api.USERVERIFICATIONENUM_PREFERRED))
	preventDuplicateDevicesDefault := helpers.BoolDefault(true)

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
			"user_verification": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userVerificationDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserVerificationEnumEnumValues), helpers.WithDefault(userVerificationDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserVerificationEnumEnumValues),
				},
			},
			"resident_key_requirement": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             residentKeyRequirementDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserVerificationEnumEnumValues), helpers.WithDefault(residentKeyRequirementDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserVerificationEnumEnumValues),
				},
			},
			// No Default, unlike the two above - it stays a plain Optional string.
			"authenticator_attachment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedAuthenticatorAttachmentEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedAuthenticatorAttachmentEnumEnumValues),
				},
			},
			"device_type_restrictions": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"max_attempts": schema.Int32Attribute{
				Optional: true,
			},
			// No MarkdownDescription: SDKv2 put the enum prose on the Elem schema, which
			// CoreConfigSchema drops for primitive lists.
			"hints": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedWebAuthnHintEnumEnumValues)),
				},
			},
			"prevent_duplicate_devices": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             preventDuplicateDevicesDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(preventDuplicateDevicesDefault.Value())),
			},
		},
	}
}

func (r *stageAuthenticatorWebAuthnResource) toRequest(ctx context.Context, data *stageAuthenticatorWebAuthnModel) (*api.AuthenticatorWebAuthnStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	deviceTypeRestrictions, d := helpers.SliceOrEmpty[string](ctx, data.DeviceTypeRestrictions)
	diags.Append(d...)

	hints, d := helpers.SliceOrEmpty[string](ctx, data.Hints)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	body := &api.AuthenticatorWebAuthnStageRequest{
		Name:                    data.Name.ValueString(),
		UserVerification:        api.UserVerificationEnum(data.UserVerification.ValueString()).Ptr(),
		ResidentKeyRequirement:  api.UserVerificationEnum(data.ResidentKeyRequirement.ValueString()).Ptr(),
		DeviceTypeRestrictions:  deviceTypeRestrictions,
		Hints:                   helpers.CastSliceString[api.WebAuthnHintEnum](hints),
		FriendlyName:            helpers.StringPtr(data.FriendlyName),
		ConfigureFlow:           *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
		MaxAttempts:             helpers.Int32Ptr(data.MaxAttempts),
		PreventDuplicateDevices: new(data.PreventDuplicateDevices.ValueBool()),
	}

	// The request model gates this field on IsSet() rather than on nil-ness, so leaving
	// the Nullable untouched is what omits it - matching SDKv2's `if x, set :=
	// d.GetOk(...)` guard. Calling Set(nil) instead would send an explicit JSON null.
	if !data.AuthenticatorAttachment.IsNull() {
		body.AuthenticatorAttachment.Set(api.AuthenticatorAttachmentEnum(data.AuthenticatorAttachment.ValueString()).Ptr())
	}

	return body, diags
}

func (r *stageAuthenticatorWebAuthnResource) fromAPI(ctx context.Context, data *stageAuthenticatorWebAuthnModel, res *api.AuthenticatorWebAuthnStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	deviceTypeRestrictions, d := helpers.MergeStringList(ctx, data.DeviceTypeRestrictions, res.DeviceTypeRestrictions)
	diags.Append(d...)
	data.DeviceTypeRestrictions = deviceTypeRestrictions

	hints := make([]string, len(res.Hints))
	for i, h := range res.Hints {
		hints[i] = string(h)
	}
	mergedHints, d := helpers.MergeStringList(ctx, data.Hints, hints)
	diags.Append(d...)
	data.Hints = mergedHints

	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
	// No Defaults on these two, so they keep the prior-aware helpers.
	data.AuthenticatorAttachment = helpers.StringOrNull(data.AuthenticatorAttachment, string(res.GetAuthenticatorAttachment()))
	data.MaxAttempts = helpers.Int32OrNull(data.MaxAttempts, res.GetMaxAttempts())
	// These do have Defaults, so they take the API value verbatim. friendly_name is the
	// one that would otherwise break import: its default is "", which is also what the
	// API returns when unset.
	data.FriendlyName = types.StringValue(res.GetFriendlyName())
	data.UserVerification = types.StringValue(string(res.GetUserVerification()))
	data.ResidentKeyRequirement = types.StringValue(string(res.GetResidentKeyRequirement()))
	data.PreventDuplicateDevices = types.BoolValue(res.GetPreventDuplicateDevices())

	return diags
}

func (r *stageAuthenticatorWebAuthnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorWebAuthnModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorWebauthnCreate(ctx).AuthenticatorWebAuthnStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorWebAuthnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorWebAuthnModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorWebauthnRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorWebAuthnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorWebAuthnModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorWebauthnUpdate(ctx, data.ID.ValueString()).AuthenticatorWebAuthnStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorWebAuthnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorWebAuthnModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorWebauthnDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
