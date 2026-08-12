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
	_ resource.Resource                = &stageIdentificationResource{}
	_ resource.ResourceWithConfigure   = &stageIdentificationResource{}
	_ resource.ResourceWithImportState = &stageIdentificationResource{}
)

func newStageIdentificationResource() resource.Resource {
	return &stageIdentificationResource{}
}

type stageIdentificationResource struct {
	resourceBase
}

type stageIdentificationModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	UserFields              types.List   `tfsdk:"user_fields"`
	PasswordStage           types.String `tfsdk:"password_stage"`
	CaptchaStage            types.String `tfsdk:"captcha_stage"`
	WebauthnStage           types.String `tfsdk:"webauthn_stage"`
	CaseInsensitiveMatching types.Bool   `tfsdk:"case_insensitive_matching"`
	ShowMatchedUser         types.Bool   `tfsdk:"show_matched_user"`
	PretendUserExists       types.Bool   `tfsdk:"pretend_user_exists"`
	ShowSourceLabels        types.Bool   `tfsdk:"show_source_labels"`
	EnableRememberMe        types.Bool   `tfsdk:"enable_remember_me"`
	EnrollmentFlow          types.String `tfsdk:"enrollment_flow"`
	RecoveryFlow            types.String `tfsdk:"recovery_flow"`
	PasswordlessFlow        types.String `tfsdk:"passwordless_flow"`
	Sources                 types.List   `tfsdk:"sources"`
}

func (r *stageIdentificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_identification"
}

func (r *stageIdentificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	showMatchedUserDefault := helpers.BoolDefault(true)
	pretendUserExistsDefault := helpers.BoolDefault(true)
	showSourceLabelsDefault := helpers.BoolDefault(false)
	enableRememberMeDefault := helpers.BoolDefault(false)

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
			// No MarkdownDescription: SDKv2 declared the enum prose on the Elem schema,
			// which CoreConfigSchema drops for primitive lists, so
			// docs/resources/stage_identification.md renders user_fields bare.
			"user_fields": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedUserFieldsEnumEnumValues)),
				},
			},
			"password_stage": schema.StringAttribute{
				Optional: true,
			},
			"captcha_stage": schema.StringAttribute{
				Optional: true,
			},
			"webauthn_stage": schema.StringAttribute{
				Optional: true,
			},
			"case_insensitive_matching": schema.BoolAttribute{
				Optional: true,
			},
			"show_matched_user": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             showMatchedUserDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(showMatchedUserDefault.Value())),
			},
			"pretend_user_exists": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             pretendUserExistsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(pretendUserExistsDefault.Value())),
			},
			"show_source_labels": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             showSourceLabelsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(showSourceLabelsDefault.Value())),
			},
			"enable_remember_me": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enableRememberMeDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enableRememberMeDefault.Value())),
			},
			"enrollment_flow": schema.StringAttribute{
				Optional: true,
			},
			"recovery_flow": schema.StringAttribute{
				Optional: true,
			},
			"passwordless_flow": schema.StringAttribute{
				Optional: true,
			},
			"sources": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *stageIdentificationResource) toRequest(ctx context.Context, data *stageIdentificationModel) (*api.IdentificationStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// SliceOrEmpty, not ElementsAs: a nil slice would be omitted from the request body,
	// so clearing either list in config would silently fail to clear it server-side.
	userFields, d := helpers.SliceOrEmpty[string](ctx, data.UserFields)
	diags.Append(d...)

	sources, d := helpers.SliceOrEmpty[string](ctx, data.Sources)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.IdentificationStageRequest{
		Name:       data.Name.ValueString(),
		UserFields: helpers.CastSliceString[api.UserFieldsEnum](userFields),
		Sources:    sources,
		// case_insensitive_matching has no Default, but SDKv2 sent
		// new(d.Get(...).(bool)) unconditionally, i.e. false when unset. ValueBool()
		// on a null types.Bool is also false, so this keeps the wire format identical
		// rather than newly omitting the field.
		CaseInsensitiveMatching: new(data.CaseInsensitiveMatching.ValueBool()),
		ShowMatchedUser:         new(data.ShowMatchedUser.ValueBool()),
		PretendUserExists:       new(data.PretendUserExists.ValueBool()),
		ShowSourceLabels:        new(data.ShowSourceLabels.ValueBool()),
		EnableRememberMe:        new(data.EnableRememberMe.ValueBool()),
		// These three stage references use StringPtrEmpty, not StringPtr: SDKv2 built
		// them with new(d.Get(...).(string)), which is a pointer to "" when unset, so
		// clearing the attribute must keep sending "" to clear it server-side.
		PasswordStage: *api.NewNullableString(helpers.StringPtrEmpty(data.PasswordStage)),
		CaptchaStage:  *api.NewNullableString(helpers.StringPtrEmpty(data.CaptchaStage)),
		WebauthnStage: *api.NewNullableString(helpers.StringPtrEmpty(data.WebauthnStage)),
		// The three flow references used GetP, which omits when unset - StringPtr.
		EnrollmentFlow:   *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		RecoveryFlow:     *api.NewNullableString(helpers.StringPtr(data.RecoveryFlow)),
		PasswordlessFlow: *api.NewNullableString(helpers.StringPtr(data.PasswordlessFlow)),
	}, diags
}

func (r *stageIdentificationResource) fromAPI(ctx context.Context, data *stageIdentificationModel, res *api.IdentificationStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	// Both lists are Optional-only, so state must match the config exactly after apply
	// or the framework raises a data-consistency error. SDKv2 merged sources but wrote
	// the API's order for user_fields directly; both merge here.
	userFields := make([]string, len(res.UserFields))
	for i, f := range res.UserFields {
		userFields[i] = string(f)
	}
	mergedUserFields, d := helpers.MergeStringList(ctx, data.UserFields, userFields)
	diags.Append(d...)
	data.UserFields = mergedUserFields

	sources, d := helpers.MergeStringList(ctx, data.Sources, res.Sources)
	diags.Append(d...)
	data.Sources = sources

	data.PasswordStage = helpers.StringPtrOrNull(res.PasswordStage.Get())
	data.CaptchaStage = helpers.StringPtrOrNull(res.CaptchaStage.Get())
	data.WebauthnStage = helpers.StringPtrOrNull(res.WebauthnStage.Get())
	data.EnrollmentFlow = helpers.StringPtrOrNull(res.EnrollmentFlow.Get())
	data.RecoveryFlow = helpers.StringPtrOrNull(res.RecoveryFlow.Get())
	data.PasswordlessFlow = helpers.StringPtrOrNull(res.PasswordlessFlow.Get())

	// case_insensitive_matching has no Default, so it keeps the prior-aware helper.
	data.CaseInsensitiveMatching = helpers.BoolOrNull(data.CaseInsensitiveMatching, res.GetCaseInsensitiveMatching())
	// The remaining four do have Defaults, so they take the API value verbatim.
	data.ShowMatchedUser = types.BoolValue(res.GetShowMatchedUser())
	data.PretendUserExists = types.BoolValue(res.GetPretendUserExists())
	data.ShowSourceLabels = types.BoolValue(res.GetShowSourceLabels())
	data.EnableRememberMe = types.BoolValue(res.GetEnableRememberMe())

	return diags
}

func (r *stageIdentificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageIdentificationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesIdentificationCreate(ctx).IdentificationStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageIdentificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageIdentificationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesIdentificationRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageIdentificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageIdentificationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesIdentificationUpdate(ctx, data.ID.ValueString()).IdentificationStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageIdentificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageIdentificationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesIdentificationDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
