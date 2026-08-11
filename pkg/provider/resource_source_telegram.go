package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &sourceTelegramResource{}
	_ resource.ResourceWithConfigure   = &sourceTelegramResource{}
	_ resource.ResourceWithImportState = &sourceTelegramResource{}
)

func newSourceTelegramResource() resource.Resource {
	return &sourceTelegramResource{}
}

type sourceTelegramResource struct {
	resourceBase
}

type sourceTelegramModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	UUID                  types.String `tfsdk:"uuid"`
	Slug                  types.String `tfsdk:"slug"`
	UserPathTemplate      types.String `tfsdk:"user_path_template"`
	AuthenticationFlow    types.String `tfsdk:"authentication_flow"`
	EnrollmentFlow        types.String `tfsdk:"enrollment_flow"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	PolicyEngineMode      types.String `tfsdk:"policy_engine_mode"`
	UserMatchingMode      types.String `tfsdk:"user_matching_mode"`
	PreAuthenticationFlow types.String `tfsdk:"pre_authentication_flow"`
	BotUsername           types.String `tfsdk:"bot_username"`
	BotToken              types.String `tfsdk:"bot_token"`
	RequestMessageAccess  types.Bool   `tfsdk:"request_message_access"`
	PropertyMappings      types.List   `tfsdk:"property_mappings"`
	PropertyMappingsGroup types.List   `tfsdk:"property_mappings_group"`
}

func (r *sourceTelegramResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_telegram"
}

func (r *sourceTelegramResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	userMatchingModeDefault := helpers.StringDefault(string(api.USERMATCHINGMODEENUM_IDENTIFIER))
	requestMessageAccessDefault := helpers.BoolDefault(false)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Directory --- ",
		Attributes: map[string]schema.Attribute{
			// H3: id is res.Slug, so no UseStateForUnknown - see resource_source_scim.go.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Required: true,
			},
			"user_path_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userPathTemplateDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(userPathTemplateDefault.Value())),
			},
			"authentication_flow": schema.StringAttribute{
				Optional: true,
			},
			"enrollment_flow": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enabledDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enabledDefault.Value())),
			},
			"policy_engine_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             policyEngineModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues), helpers.WithDefault(policyEngineModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedPolicyEngineModeEnumValues),
				},
			},
			"user_matching_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userMatchingModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues), helpers.WithDefault(userMatchingModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserMatchingModeEnumEnumValues),
				},
			},
			"pre_authentication_flow": schema.StringAttribute{
				Required: true,
			},
			"bot_username": schema.StringAttribute{
				Required: true,
			},
			// One of the 18 H4 attributes: TelegramSource has no bot_token field at all, so
			// the API cannot return it. Note it is deliberately *not* marked Sensitive -
			// SDKv2 did not mark it either, and doing so here would change the doc page.
			"bot_token": schema.StringAttribute{
				Required: true,
			},
			"request_message_access": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             requestMessageAccessDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(requestMessageAccessDefault.Value())),
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"property_mappings_group": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *sourceTelegramResource) toRequest(ctx context.Context, data *sourceTelegramModel) (*api.TelegramSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.TelegramSourceRequest{
		Name:                  data.Name.ValueString(),
		Slug:                  data.Slug.ValueString(),
		Enabled:               new(data.Enabled.ValueBool()),
		UserPathTemplate:      new(data.UserPathTemplate.ValueString()),
		PolicyEngineMode:      api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		UserMatchingMode:      api.UserMatchingModeEnum(data.UserMatchingMode.ValueString()).Ptr(),
		AuthenticationFlow:    *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		EnrollmentFlow:        *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		PreAuthenticationFlow: data.PreAuthenticationFlow.ValueString(),
		UserPropertyMappings:  propertyMappings,
		GroupPropertyMappings: propertyMappingsGroup,
		BotUsername:           data.BotUsername.ValueString(),
		BotToken:              data.BotToken.ValueString(),
		// Another GetP[bool] casualty, same as authentik_policy_password's three checks:
		// SDKv2 omitted request_message_access whenever it was false, and since the API
		// leaves absent fields unchanged on PUT it could not be turned back off. It has a
		// Default, so it is never null here and the concrete value is always sent.
		RequestMessageAccess: new(data.RequestMessageAccess.ValueBool()),
	}, diags
}

// fromAPI deliberately never assigns BotToken: this is one of the 18 H4 resources, and the
// TelegramSource response type has no bot_token field at all. Callers seed data from
// req.Plan or req.State first, so the omission preserves the configured secret - see
// stageCaptchaResource.fromAPI for the full explanation. bot_token is Required, so nulling
// it would fail the apply rather than merely losing the value.
func (r *sourceTelegramResource) fromAPI(ctx context.Context, data *sourceTelegramModel, res *api.TelegramSource) diag.Diagnostics {
	var diags diag.Diagnostics

	// H3: id is the slug, not the Pk.
	data.ID = types.StringValue(res.Slug)
	data.Name = types.StringValue(res.Name)
	data.Slug = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.UserPropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	propertyMappingsGroup, d := helpers.MergeStringList(ctx, data.PropertyMappingsGroup, res.GroupPropertyMappings)
	diags.Append(d...)
	data.PropertyMappingsGroup = propertyMappingsGroup

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.EnrollmentFlow = helpers.StringPtrOrNull(res.EnrollmentFlow.Get())
	// Required, so never null.
	data.PreAuthenticationFlow = types.StringValue(res.PreAuthenticationFlow)
	data.BotUsername = types.StringValue(res.BotUsername)
	// Everything below has a Default and takes the API value verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.PolicyEngineMode = types.StringValue(string(res.GetPolicyEngineMode()))
	data.UserMatchingMode = types.StringValue(string(res.GetUserMatchingMode()))
	data.RequestMessageAccess = types.BoolValue(res.GetRequestMessageAccess())

	return diags
}

func (r *sourceTelegramResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceTelegramModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesTelegramCreate(ctx).TelegramSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceTelegramResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries bot_token across a Read - see fromAPI.
	var data sourceTelegramModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesTelegramRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *sourceTelegramResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceTelegramModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// H3/discovery #4: the plan's id is unknown when the slug changes.
	var priorID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &priorID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesTelegramUpdate(ctx, priorID.ValueString()).TelegramSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceTelegramResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceTelegramModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesTelegramDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
