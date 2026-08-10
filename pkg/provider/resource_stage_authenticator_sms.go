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
	_ resource.Resource                = &stageAuthenticatorSmsResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorSmsResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorSmsResource{}
)

func newStageAuthenticatorSmsResource() resource.Resource {
	return &stageAuthenticatorSmsResource{}
}

type stageAuthenticatorSmsResource struct {
	resourceBase
}

type stageAuthenticatorSmsModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	FriendlyName  types.String `tfsdk:"friendly_name"`
	ConfigureFlow types.String `tfsdk:"configure_flow"`
	SMSProvider   types.String `tfsdk:"sms_provider"`
	FromNumber    types.String `tfsdk:"from_number"`
	AccountSid    types.String `tfsdk:"account_sid"`
	Auth          types.String `tfsdk:"auth"`
	AuthType      types.String `tfsdk:"auth_type"`
	AuthPassword  types.String `tfsdk:"auth_password"`
	Mapping       types.String `tfsdk:"mapping"`
	VerifyOnly    types.Bool   `tfsdk:"verify_only"`
}

func (r *stageAuthenticatorSmsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_sms"
}

func (r *stageAuthenticatorSmsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")
	smsProviderDefault := helpers.StringDefault(string(api.PROVIDERENUM_TWILIO))
	authTypeDefault := helpers.StringDefault(string(api.AUTHTYPEENUM_BASIC))
	verifyOnlyDefault := helpers.BoolDefault(false)

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
			// Named sms_provider in Terraform but "provider" on the wire - the attribute
			// could not be called provider inside a provider block.
			"sms_provider": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             smsProviderDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedProviderEnumEnumValues), helpers.WithDefault(smsProviderDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedProviderEnumEnumValues),
				},
			},
			"from_number": schema.StringAttribute{
				Required: true,
			},
			"account_sid": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"auth": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"auth_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             authTypeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedAuthTypeEnumEnumValues), helpers.WithDefault(authTypeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedAuthTypeEnumEnumValues),
				},
			},
			"auth_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"mapping": schema.StringAttribute{
				Optional: true,
			},
			"verify_only": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             verifyOnlyDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(verifyOnlyDefault.Value())),
			},
		},
	}
}

func (r *stageAuthenticatorSmsResource) toRequest(data *stageAuthenticatorSmsModel) *api.AuthenticatorSMSStageRequest {
	return &api.AuthenticatorSMSStageRequest{
		Name:         data.Name.ValueString(),
		Provider:     api.ProviderEnum(data.SMSProvider.ValueString()),
		FromNumber:   data.FromNumber.ValueString(),
		AccountSid:   data.AccountSid.ValueString(),
		Auth:         data.Auth.ValueString(),
		AuthType:     api.AuthTypeEnum(data.AuthType.ValueString()).Ptr(),
		FriendlyName: helpers.StringPtr(data.FriendlyName),
		AuthPassword: helpers.StringPtr(data.AuthPassword),
		// SDKv2 built this with GetP[bool], which returns nil whenever the value is
		// false, so verify_only was omitted from the request body and an existing
		// verify_only = true could never be turned back off (the API keeps absent fields
		// unchanged on PUT, Read returned true, and the plan showed a permanent diff).
		// The attribute has a Default, so it is never null here and the concrete value is
		// always sent.
		VerifyOnly:    new(data.VerifyOnly.ValueBool()),
		ConfigureFlow: *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
		Mapping:       *api.NewNullableString(helpers.StringPtr(data.Mapping)),
	}
}

func (r *stageAuthenticatorSmsResource) fromAPI(data *stageAuthenticatorSmsModel, res *api.AuthenticatorSMSStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.FromNumber = types.StringValue(res.FromNumber)
	data.AccountSid = types.StringValue(res.AccountSid)
	data.Auth = types.StringValue(res.Auth)
	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
	data.Mapping = helpers.StringPtrOrNull(res.Mapping.Get())
	// auth_password has no Default, so it keeps the prior-aware helper.
	data.AuthPassword = helpers.StringOrNull(data.AuthPassword, res.GetAuthPassword())
	// The rest have Defaults and take the API value verbatim.
	data.FriendlyName = types.StringValue(res.GetFriendlyName())
	data.SMSProvider = types.StringValue(string(res.Provider))
	data.AuthType = types.StringValue(string(res.GetAuthType()))
	data.VerifyOnly = types.BoolValue(res.GetVerifyOnly())
}

func (r *stageAuthenticatorSmsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorSmsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorSmsCreate(ctx).AuthenticatorSMSStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorSmsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorSmsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorSmsRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorSmsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorSmsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorSmsUpdate(ctx, data.ID.ValueString()).AuthenticatorSMSStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorSmsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorSmsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorSmsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
