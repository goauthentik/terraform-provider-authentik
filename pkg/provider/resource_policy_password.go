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
	_ resource.Resource                = &policyPasswordResource{}
	_ resource.ResourceWithConfigure   = &policyPasswordResource{}
	_ resource.ResourceWithImportState = &policyPasswordResource{}
)

func newPolicyPasswordResource() resource.Resource {
	return &policyPasswordResource{}
}

type policyPasswordResource struct {
	resourceBase
}

type policyPasswordModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ExecutionLogging     types.Bool   `tfsdk:"execution_logging"`
	PasswordField        types.String `tfsdk:"password_field"`
	CheckStaticRules     types.Bool   `tfsdk:"check_static_rules"`
	CheckHaveIBeenPwned  types.Bool   `tfsdk:"check_have_i_been_pwned"`
	CheckZxcvbn          types.Bool   `tfsdk:"check_zxcvbn"`
	ErrorMessage         types.String `tfsdk:"error_message"`
	AmountUppercase      types.Int32  `tfsdk:"amount_uppercase"`
	AmountLowercase      types.Int32  `tfsdk:"amount_lowercase"`
	AmountSymbols        types.Int32  `tfsdk:"amount_symbols"`
	AmountDigits         types.Int32  `tfsdk:"amount_digits"`
	LengthMin            types.Int32  `tfsdk:"length_min"`
	SymbolCharset        types.String `tfsdk:"symbol_charset"`
	HibpAllowedCount     types.Int32  `tfsdk:"hibp_allowed_count"`
	ZxcvbnScoreThreshold types.Int32  `tfsdk:"zxcvbn_score_threshold"`
}

func (r *policyPasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_password"
}

func (r *policyPasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)
	passwordFieldDefault := helpers.StringDefault("password")
	checkStaticRulesDefault := helpers.BoolDefault(true)
	checkHaveIBeenPwnedDefault := helpers.BoolDefault(false)
	checkZxcvbnDefault := helpers.BoolDefault(false)
	// Kept exactly as SDKv2 spelled it, escapes included, or the doc page changes.
	symbolCharsetDefault := helpers.StringDefault("!\\\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	hibpAllowedCountDefault := helpers.Int32Default(1)
	zxcvbnScoreThresholdDefault := helpers.Int32Default(2)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- ",
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
			"execution_logging": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             executionLoggingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(executionLoggingDefault.Value())),
			},
			"password_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             passwordFieldDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(passwordFieldDefault.Value())),
			},
			"check_static_rules": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             checkStaticRulesDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(checkStaticRulesDefault.Value())),
			},
			"check_have_i_been_pwned": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             checkHaveIBeenPwnedDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(checkHaveIBeenPwnedDefault.Value())),
			},
			"check_zxcvbn": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             checkZxcvbnDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(checkZxcvbnDefault.Value())),
			},
			"error_message": schema.StringAttribute{
				Required: true,
			},
			// The five amount_*/length_min attributes are Optional with no Default, so they
			// are the ones keeping the prior-aware read helpers here.
			"amount_uppercase": schema.Int32Attribute{
				Optional: true,
			},
			"amount_lowercase": schema.Int32Attribute{
				Optional: true,
			},
			"amount_symbols": schema.Int32Attribute{
				Optional: true,
			},
			"amount_digits": schema.Int32Attribute{
				Optional: true,
			},
			"length_min": schema.Int32Attribute{
				Optional: true,
			},
			"symbol_charset": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             symbolCharsetDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(symbolCharsetDefault.Value())),
			},
			"hibp_allowed_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             hibpAllowedCountDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(hibpAllowedCountDefault.Value())),
			},
			"zxcvbn_score_threshold": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             zxcvbnScoreThresholdDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(zxcvbnScoreThresholdDefault.Value())),
			},
		},
	}
}

func (r *policyPasswordResource) toRequest(data *policyPasswordModel) *api.PasswordPolicyRequest {
	return &api.PasswordPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
		PasswordField:    helpers.StringPtr(data.PasswordField),
		// All three of these were built with GetP[bool] in SDKv2, which returns nil for any
		// false value, so none of them was ever sent as false - and since the API leaves
		// absent fields unchanged on PUT, check_static_rules in particular (Default true)
		// could not be switched off at all. They have Defaults, so they are never null here
		// and the concrete value is always sent. Same defect as
		// authentik_stage_authenticator_sms.verify_only in batch 1; the no-Default variant
		// is in resource_policy_geoip.go, which needs helpers.BoolPtr instead.
		CheckStaticRules:    new(data.CheckStaticRules.ValueBool()),
		CheckHaveIBeenPwned: new(data.CheckHaveIBeenPwned.ValueBool()),
		CheckZxcvbn:         new(data.CheckZxcvbn.ValueBool()),
		SymbolCharset:       helpers.StringPtr(data.SymbolCharset),
		// error_message is Required in Terraform but optional on the wire, and SDKv2 sourced
		// it via GetP too - so an intentionally empty message was dropped from the request.
		// Required means it is never null here, so "" is now sent as "".
		ErrorMessage:         helpers.StringPtr(data.ErrorMessage),
		AmountUppercase:      helpers.Int32Ptr(data.AmountUppercase),
		AmountLowercase:      helpers.Int32Ptr(data.AmountLowercase),
		AmountSymbols:        helpers.Int32Ptr(data.AmountSymbols),
		AmountDigits:         helpers.Int32Ptr(data.AmountDigits),
		LengthMin:            helpers.Int32Ptr(data.LengthMin),
		HibpAllowedCount:     helpers.Int32Ptr(data.HibpAllowedCount),
		ZxcvbnScoreThreshold: helpers.Int32Ptr(data.ZxcvbnScoreThreshold),
	}
}

func (r *policyPasswordResource) fromAPI(data *policyPasswordModel, res *api.PasswordPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	// error_message is Required, so it is never null and takes the value directly.
	data.ErrorMessage = types.StringValue(res.GetErrorMessage())
	// No Default on these five, so they keep the prior-aware helper: a never-configured
	// rule stays null rather than becoming 0.
	data.AmountUppercase = helpers.Int32OrNull(data.AmountUppercase, res.GetAmountUppercase())
	data.AmountLowercase = helpers.Int32OrNull(data.AmountLowercase, res.GetAmountLowercase())
	data.AmountSymbols = helpers.Int32OrNull(data.AmountSymbols, res.GetAmountSymbols())
	data.AmountDigits = helpers.Int32OrNull(data.AmountDigits, res.GetAmountDigits())
	data.LengthMin = helpers.Int32OrNull(data.LengthMin, res.GetLengthMin())
	// The rest all have Defaults and take the API value verbatim.
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	data.PasswordField = types.StringValue(res.GetPasswordField())
	data.CheckStaticRules = types.BoolValue(res.GetCheckStaticRules())
	data.CheckHaveIBeenPwned = types.BoolValue(res.GetCheckHaveIBeenPwned())
	data.CheckZxcvbn = types.BoolValue(res.GetCheckZxcvbn())
	data.SymbolCharset = types.StringValue(res.GetSymbolCharset())
	data.HibpAllowedCount = types.Int32Value(res.GetHibpAllowedCount())
	data.ZxcvbnScoreThreshold = types.Int32Value(res.GetZxcvbnScoreThreshold())
}

func (r *policyPasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyPasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesPasswordCreate(ctx).PasswordPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyPasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyPasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesPasswordRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyPasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyPasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesPasswordUpdate(ctx, data.ID.ValueString()).PasswordPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyPasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyPasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesPasswordDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
