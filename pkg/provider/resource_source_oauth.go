package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	_ resource.Resource                = &sourceOAuthResource{}
	_ resource.ResourceWithConfigure   = &sourceOAuthResource{}
	_ resource.ResourceWithImportState = &sourceOAuthResource{}
)

func newSourceOAuthResource() resource.Resource {
	return &sourceOAuthResource{}
}

type sourceOAuthResource struct {
	resourceBase
}

type sourceOAuthModel struct {
	ID                          types.String         `tfsdk:"id"`
	Name                        types.String         `tfsdk:"name"`
	UUID                        types.String         `tfsdk:"uuid"`
	Slug                        types.String         `tfsdk:"slug"`
	UserPathTemplate            types.String         `tfsdk:"user_path_template"`
	AuthenticationFlow          types.String         `tfsdk:"authentication_flow"`
	EnrollmentFlow              types.String         `tfsdk:"enrollment_flow"`
	Enabled                     types.Bool           `tfsdk:"enabled"`
	Promoted                    types.Bool           `tfsdk:"promoted"`
	AuthorizationCodeAuthMethod types.String         `tfsdk:"authorization_code_auth_method"`
	PolicyEngineMode            types.String         `tfsdk:"policy_engine_mode"`
	UserMatchingMode            types.String         `tfsdk:"user_matching_mode"`
	GroupMatchingMode           types.String         `tfsdk:"group_matching_mode"`
	ProviderType                types.String         `tfsdk:"provider_type"`
	RequestTokenURL             types.String         `tfsdk:"request_token_url"`
	AuthorizationURL            types.String         `tfsdk:"authorization_url"`
	AccessTokenURL              types.String         `tfsdk:"access_token_url"`
	ProfileURL                  types.String         `tfsdk:"profile_url"`
	OIDCWellKnownURL            types.String         `tfsdk:"oidc_well_known_url"`
	OIDCJWKSURL                 types.String         `tfsdk:"oidc_jwks_url"`
	OIDCJWKS                    jsontypes.Normalized `tfsdk:"oidc_jwks"`
	PKCE                        types.String         `tfsdk:"pkce"`
	AdditionalScopes            types.String         `tfsdk:"additional_scopes"`
	ConsumerKey                 types.String         `tfsdk:"consumer_key"`
	ConsumerSecret              types.String         `tfsdk:"consumer_secret"`
	CallbackURI                 types.String         `tfsdk:"callback_uri"`
	PropertyMappings            types.List           `tfsdk:"property_mappings"`
	PropertyMappingsGroup       types.List           `tfsdk:"property_mappings_group"`
}

func (r *sourceOAuthResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_oauth"
}

func (r *sourceOAuthResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	promotedDefault := helpers.BoolDefault(false)
	authorizationCodeAuthMethodDefault := helpers.StringDefault(string(api.AUTHORIZATIONCODEAUTHMETHODENUM_BASIC_AUTH))
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	userMatchingModeDefault := helpers.StringDefault(string(api.USERMATCHINGMODEENUM_IDENTIFIER))
	groupMatchingModeDefault := helpers.StringDefault(string(api.GROUPMATCHINGMODEENUM_IDENTIFIER))
	pkceDefault := helpers.StringDefault(string(api.PKCEMETHODENUM_NONE))

	const manualURLDescription = "Manually configure OAuth2 URLs when `oidc_well_known_url` is not set."

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
			"promoted": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             promotedDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(promotedDefault.Value())),
			},
			"authorization_code_auth_method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             authorizationCodeAuthMethodDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues), helpers.WithDefault(authorizationCodeAuthMethodDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues),
				},
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
			"group_matching_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             groupMatchingModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues), helpers.WithDefault(groupMatchingModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedGroupMatchingModeEnumEnumValues),
				},
			},
			"provider_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedProviderTypeEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedProviderTypeEnumEnumValues),
				},
			},
			"request_token_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: manualURLDescription,
			},
			"authorization_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: manualURLDescription,
			},
			"access_token_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Only required for OAuth1.",
			},
			"profile_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: manualURLDescription,
			},
			"oidc_well_known_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Automatically configure source from OIDC well-known endpoint. URL is taken as is, and should end with `.well-known/openid-configuration`.",
			},
			"oidc_jwks_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Automatically configure JWKS if not specified by `oidc_well_known_url`.",
			},
			// Optional+Computed with no Default, and the server fills it in from
			// oidc_well_known_url - i.e. it depends on a *sibling* attribute, so per the
			// UseStateForUnknown rule it must not pin prior state. The cost is
			// "(known after apply)" on plans that change the well-known URL.
			"oidc_jwks": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				MarkdownDescription: helpers.Desc("Manually configure JWKS keys for use with machine-to-machine authentication. "+helpers.JSONDescription, helpers.Generated()),
			},
			"pkce": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             pkceDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedPKCEMethodEnumEnumValues), helpers.WithDefault(pkceDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedPKCEMethodEnumEnumValues),
				},
			},
			"additional_scopes": schema.StringAttribute{
				Optional: true,
			},
			"consumer_key": schema.StringAttribute{
				Required: true,
			},
			// One of the 18 H4 attributes: the API never returns consumer_secret.
			"consumer_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			// Computed-only and derived from the slug, so no UseStateForUnknown either.
			"callback_uri": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
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

func (r *sourceOAuthResource) toRequest(ctx context.Context, data *sourceOAuthModel) (*api.OAuthSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	var oidcJWKS map[string]any
	diags.Append(data.OIDCJWKS.Unmarshal(&oidcJWKS)...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.OAuthSourceRequest{
		Name:                        data.Name.ValueString(),
		Slug:                        data.Slug.ValueString(),
		Enabled:                     new(data.Enabled.ValueBool()),
		Promoted:                    new(data.Promoted.ValueBool()),
		UserPathTemplate:            new(data.UserPathTemplate.ValueString()),
		ProviderType:                api.ProviderTypeEnum(data.ProviderType.ValueString()),
		ConsumerKey:                 data.ConsumerKey.ValueString(),
		ConsumerSecret:              data.ConsumerSecret.ValueString(),
		AuthorizationCodeAuthMethod: api.AuthorizationCodeAuthMethodEnum(data.AuthorizationCodeAuthMethod.ValueString()).Ptr(),
		PolicyEngineMode:            api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		UserMatchingMode:            api.UserMatchingModeEnum(data.UserMatchingMode.ValueString()).Ptr(),
		GroupMatchingMode:           api.GroupMatchingModeEnum(data.GroupMatchingMode.ValueString()).Ptr(),
		AuthenticationFlow:          *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		EnrollmentFlow:              *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		RequestTokenUrl:             *api.NewNullableString(helpers.StringPtr(data.RequestTokenURL)),
		AuthorizationUrl:            *api.NewNullableString(helpers.StringPtr(data.AuthorizationURL)),
		AccessTokenUrl:              *api.NewNullableString(helpers.StringPtr(data.AccessTokenURL)),
		ProfileUrl:                  *api.NewNullableString(helpers.StringPtr(data.ProfileURL)),
		AdditionalScopes:            helpers.StringPtr(data.AdditionalScopes),
		OidcWellKnownUrl:            helpers.StringPtr(data.OIDCWellKnownURL),
		OidcJwksUrl:                 helpers.StringPtr(data.OIDCJWKSURL),
		OidcJwks:                    oidcJWKS,
		Pkce:                        api.PKCEMethodEnum(data.PKCE.ValueString()).Ptr(),
		UserPropertyMappings:        propertyMappings,
		GroupPropertyMappings:       propertyMappingsGroup,
	}, diags
}

// fromAPI deliberately never assigns ConsumerSecret: this is one of the 18 H4 resources and
// the API does not return it. Callers seed data from req.Plan or req.State first, so the
// omission preserves the configured secret - see stageCaptchaResource.fromAPI.
func (r *sourceOAuthResource) fromAPI(ctx context.Context, data *sourceOAuthModel, res *api.OAuthSource) diag.Diagnostics {
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
	data.RequestTokenURL = helpers.StringPtrOrNull(res.RequestTokenUrl.Get())
	data.AuthorizationURL = helpers.StringPtrOrNull(res.AuthorizationUrl.Get())
	data.AccessTokenURL = helpers.StringPtrOrNull(res.AccessTokenUrl.Get())
	data.ProfileURL = helpers.StringPtrOrNull(res.ProfileUrl.Get())
	// Plain *string with no Default, so prior-aware.
	data.AdditionalScopes = helpers.StringOrNull(data.AdditionalScopes, res.GetAdditionalScopes())
	data.OIDCWellKnownURL = helpers.StringOrNull(data.OIDCWellKnownURL, res.GetOidcWellKnownUrl())
	data.OIDCJWKSURL = helpers.StringOrNull(data.OIDCJWKSURL, res.GetOidcJwksUrl())
	// Required, so never null.
	data.ProviderType = types.StringValue(string(res.ProviderType))
	data.ConsumerKey = types.StringValue(res.ConsumerKey)
	// Computed-only.
	data.CallbackURI = types.StringValue(res.CallbackUrl)
	// Have Defaults, so verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.Promoted = types.BoolValue(res.GetPromoted())
	data.AuthorizationCodeAuthMethod = types.StringValue(string(res.GetAuthorizationCodeAuthMethod()))
	data.PolicyEngineMode = types.StringValue(string(res.GetPolicyEngineMode()))
	data.UserMatchingMode = types.StringValue(string(res.GetUserMatchingMode()))
	data.GroupMatchingMode = types.StringValue(string(res.GetGroupMatchingMode()))
	data.PKCE = types.StringValue(string(res.GetPkce()))

	jwksBytes, err := json.Marshal(res.OidcJwks)
	if err != nil {
		diags.AddError("Failed to encode OAuth source oidc_jwks", err.Error())
		return diags
	}
	data.OIDCJWKS = jsontypes.NewNormalizedValue(string(jwksBytes))

	return diags
}

func (r *sourceOAuthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceOAuthModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesOauthCreate(ctx).OAuthSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceOAuthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries consumer_secret across a Read - see fromAPI.
	var data sourceOAuthModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesOauthRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *sourceOAuthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceOAuthModel
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

	res, hr, err := r.client.SourcesAPI.SourcesOauthUpdate(ctx, priorID.ValueString()).OAuthSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceOAuthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceOAuthModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesOauthDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
