package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &providerOAuth2Resource{}
	_ resource.ResourceWithConfigure   = &providerOAuth2Resource{}
	_ resource.ResourceWithImportState = &providerOAuth2Resource{}
)

// redirectURIElementType is the element type of allowed_redirect_uris.
//
// SDKv2 declared it as TypeList of a bare TypeMap, and CoreConfigSchema expands a TypeMap
// with no Elem to map(string) - so the protocol type is list(map(string)) and the doc page
// says "List of Map of String". Keeping it as ListAttribute{MapType{StringType}} preserves
// config and state byte-for-byte. Converting it to a ListNestedAttribute would emit a
// protocol NestedType instead, changing the state layout and requiring a state upgrader, so
// per migration-plan.md that is deliberately left as a separate follow-up.
var redirectURIElementType = types.MapType{ElemType: types.StringType}

func newProviderOAuth2Resource() resource.Resource {
	return &providerOAuth2Resource{}
}

type providerOAuth2Resource struct {
	resourceBase
}

type providerOAuth2Model struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	AuthenticationFlow     types.String `tfsdk:"authentication_flow"`
	AuthorizationFlow      types.String `tfsdk:"authorization_flow"`
	InvalidationFlow       types.String `tfsdk:"invalidation_flow"`
	PropertyMappings       types.List   `tfsdk:"property_mappings"`
	ClientType             types.String `tfsdk:"client_type"`
	GrantTypes             types.List   `tfsdk:"grant_types"`
	ClientID               types.String `tfsdk:"client_id"`
	ClientSecret           types.String `tfsdk:"client_secret"`
	AccessCodeValidity     types.String `tfsdk:"access_code_validity"`
	AccessTokenValidity    types.String `tfsdk:"access_token_validity"`
	RefreshTokenValidity   types.String `tfsdk:"refresh_token_validity"`
	RefreshTokenThreshold  types.String `tfsdk:"refresh_token_threshold"`
	IncludeClaimsInIDToken types.Bool   `tfsdk:"include_claims_in_id_token"`
	SigningKey             types.String `tfsdk:"signing_key"`
	EncryptionKey          types.String `tfsdk:"encryption_key"`
	AllowedRedirectURIs    types.List   `tfsdk:"allowed_redirect_uris"`
	LogoutMethod           types.String `tfsdk:"logout_method"`
	LogoutURI              types.String `tfsdk:"logout_uri"`
	SubMode                types.String `tfsdk:"sub_mode"`
	IssuerMode             types.String `tfsdk:"issuer_mode"`
	JWKSSources            types.List   `tfsdk:"jwks_sources"`
	JWTFederationSources   types.List   `tfsdk:"jwt_federation_sources"`
	JWTFederationProviders types.List   `tfsdk:"jwt_federation_providers"`
}

func (r *providerOAuth2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_oauth2"
}

func (r *providerOAuth2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	clientTypeDefault := helpers.StringDefault(string(api.CLIENTTYPEENUM_CONFIDENTIAL))
	accessCodeValidityDefault := helpers.StringDefault("minutes=1")
	accessTokenValidityDefault := helpers.StringDefault("minutes=10")
	refreshTokenValidityDefault := helpers.StringDefault("days=30")
	refreshTokenThresholdDefault := helpers.StringDefault("seconds=0")
	includeClaimsInIDTokenDefault := helpers.BoolDefault(true)
	logoutMethodDefault := helpers.StringDefault(string(api.OAUTH2PROVIDERLOGOUTMETHODENUM_BACKCHANNEL))
	subModeDefault := helpers.StringDefault(string(api.SUBMODEENUM_HASHED_USER_ID))
	issuerModeDefault := helpers.StringDefault(string(api.ISSUERMODEENUM_PER_PROVIDER))

	resp.Schema = schema.Schema{
		MarkdownDescription: "Applications --- ",
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
			"authentication_flow": schema.StringAttribute{
				Optional: true,
			},
			"authorization_flow": schema.StringAttribute{
				Required: true,
			},
			"invalidation_flow": schema.StringAttribute{
				Required: true,
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"client_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             clientTypeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedClientTypeEnumEnumValues), helpers.WithDefault(clientTypeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedClientTypeEnumEnumValues),
				},
			},
			// grant_types is the canonical sibling-derived attribute: the server computes it
			// from client_type, so it must NOT carry UseStateForUnknown. Pinning prior state
			// would make the plan assert the stale list while the API returns a new one, and
			// the apply would fail. The cost is "(known after apply)" whenever anything it
			// depends on changes, which migration-plan.md accepts explicitly.
			//
			// No MarkdownDescription for the enum: SDKv2 put that prose on the Elem schema,
			// which CoreConfigSchema drops for primitive lists, so the page shows only
			// "Generated." here.
			"grant_types": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedGrantTypesEnumEnumValues)),
				},
			},
			"client_id": schema.StringAttribute{
				Required: true,
			},
			// Generated once at creation and stable thereafter - it does not depend on a
			// sibling - so unlike grant_types this one is safe to pin.
			"client_secret": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_code_validity": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             accessCodeValidityDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(accessCodeValidityDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"access_token_validity": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             accessTokenValidityDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(accessTokenValidityDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"refresh_token_validity": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             refreshTokenValidityDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(refreshTokenValidityDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"refresh_token_threshold": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             refreshTokenThresholdDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(refreshTokenThresholdDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"include_claims_in_id_token": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             includeClaimsInIDTokenDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(includeClaimsInIDTokenDefault.Value())),
			},
			"signing_key": schema.StringAttribute{
				Optional: true,
			},
			"encryption_key": schema.StringAttribute{
				Optional: true,
			},
			"allowed_redirect_uris": schema.ListAttribute{
				ElementType: redirectURIElementType,
				Optional:    true,
			},
			"logout_method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             logoutMethodDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedOAuth2ProviderLogoutMethodEnumEnumValues), helpers.WithDefault(logoutMethodDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedOAuth2ProviderLogoutMethodEnumEnumValues),
				},
			},
			"logout_uri": schema.StringAttribute{
				Optional: true,
			},
			"sub_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             subModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedSubModeEnumEnumValues), helpers.WithDefault(subModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedSubModeEnumEnumValues),
				},
			},
			"issuer_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             issuerModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedIssuerModeEnumEnumValues), helpers.WithDefault(issuerModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedIssuerModeEnumEnumValues),
				},
			},
			// jwks_sources is inert: SDKv2 declared it but never sent it in a request and
			// never read it back, so it only ever holds whatever config put there. Kept
			// exactly that way - see fromAPI.
			"jwks_sources": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Deprecated. Use `jwt_federation_sources` instead.",
			},
			"jwt_federation_sources": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "JWTs issued by keys configured in any of the selected sources can be used to authenticate on behalf of this provider.",
			},
			"jwt_federation_providers": schema.ListAttribute{
				ElementType:         types.Int32Type,
				Optional:            true,
				MarkdownDescription: "JWTs issued by any of the configured providers can be used to authenticate on behalf of this provider.",
			},
		},
	}
}

// customRedirectURI mirrors pkg/sdkprovider's CustomRedirectURI: a comparable struct so
// helpers.ListConsistentMerge can match entries by value and keep the configured ordering.
type customRedirectURI struct {
	MatchingMode    api.MatchingModeEnum
	URL             string
	RedirectURIType api.RedirectURITypeEnum
}

// redirectURIsFromList converts allowed_redirect_uris into comparable structs. Each element
// is a map(string) whose redirect_uri_type key is optional, exactly as SDKv2 wrote it.
func redirectURIsFromList(ctx context.Context, list types.List) ([]customRedirectURI, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := []customRedirectURI{}

	if list.IsNull() || list.IsUnknown() {
		return out, diags
	}

	var raw []map[string]string
	diags.Append(list.ElementsAs(ctx, &raw, false)...)
	if diags.HasError() {
		return out, diags
	}

	for _, entry := range raw {
		out = append(out, customRedirectURI{
			MatchingMode:    api.MatchingModeEnum(entry["matching_mode"]),
			URL:             entry["url"],
			RedirectURIType: api.RedirectURITypeEnum(entry["redirect_uri_type"]),
		})
	}
	return out, diags
}

// redirectURIsToList is the inverse. It omits redirect_uri_type when empty, matching SDKv2's
// redirectURIsToList - so an element map has two keys or three depending on whether the type
// is set, which map(string) allows.
func redirectURIsToList(uris []customRedirectURI) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	elems := make([]attr.Value, 0, len(uris))
	for _, u := range uris {
		entry := map[string]attr.Value{
			"matching_mode": types.StringValue(string(u.MatchingMode)),
			"url":           types.StringValue(u.URL),
		}
		if u.RedirectURIType != "" {
			entry["redirect_uri_type"] = types.StringValue(string(u.RedirectURIType))
		}
		m, d := types.MapValue(types.StringType, entry)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(redirectURIElementType), diags
		}
		elems = append(elems, m)
	}

	list, d := types.ListValue(redirectURIElementType, elems)
	diags.Append(d...)
	return list, diags
}

func redirectURIsFromAPI(raw []api.RedirectURI) []customRedirectURI {
	out := make([]customRedirectURI, len(raw))
	for i, u := range raw {
		c := customRedirectURI{MatchingMode: u.MatchingMode, URL: u.Url}
		if u.RedirectUriType != nil {
			c.RedirectURIType = *u.RedirectUriType
		}
		out[i] = c
	}
	return out
}

func (r *providerOAuth2Resource) toRequest(ctx context.Context, data *providerOAuth2Model) (*api.OAuth2ProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	jwtFederationSources, d := helpers.SliceOrEmpty[string](ctx, data.JWTFederationSources)
	diags.Append(d...)

	jwtFederationProviders, d := helpers.SliceOrEmpty[int32](ctx, data.JWTFederationProviders)
	diags.Append(d...)

	redirectURIs, d := redirectURIsFromList(ctx, data.AllowedRedirectURIs)
	diags.Append(d...)

	grantTypes, d := helpers.SliceOrEmpty[string](ctx, data.GrantTypes)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	redirectURIRequests := make([]api.RedirectURIRequest, 0, len(redirectURIs))
	for _, u := range redirectURIs {
		req := api.RedirectURIRequest{MatchingMode: u.MatchingMode, Url: u.URL}
		if u.RedirectURIType != "" {
			req.RedirectUriType = u.RedirectURIType.Ptr()
		}
		redirectURIRequests = append(redirectURIRequests, req)
	}

	body := &api.OAuth2ProviderRequest{
		Name:                   data.Name.ValueString(),
		AuthorizationFlow:      data.AuthorizationFlow.ValueString(),
		AuthenticationFlow:     *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		InvalidationFlow:       data.InvalidationFlow.ValueString(),
		AccessCodeValidity:     new(data.AccessCodeValidity.ValueString()),
		AccessTokenValidity:    new(data.AccessTokenValidity.ValueString()),
		RefreshTokenValidity:   new(data.RefreshTokenValidity.ValueString()),
		RefreshTokenThreshold:  helpers.StringPtr(data.RefreshTokenThreshold),
		IncludeClaimsInIdToken: new(data.IncludeClaimsInIDToken.ValueBool()),
		ClientId:               new(data.ClientID.ValueString()),
		ClientSecret:           helpers.StringPtr(data.ClientSecret),
		IssuerMode:             api.IssuerModeEnum(data.IssuerMode.ValueString()).Ptr(),
		SubMode:                api.SubModeEnum(data.SubMode.ValueString()).Ptr(),
		ClientType:             api.ClientTypeEnum(data.ClientType.ValueString()).Ptr(),
		PropertyMappings:       propertyMappings,
		JwtFederationSources:   jwtFederationSources,
		JwtFederationProviders: jwtFederationProviders,
		LogoutMethod:           api.OAuth2ProviderLogoutMethodEnum(data.LogoutMethod.ValueString()).Ptr(),
		LogoutUri:              helpers.StringPtr(data.LogoutURI),
		SigningKey:             *api.NewNullableString(helpers.StringPtr(data.SigningKey)),
		EncryptionKey:          *api.NewNullableString(helpers.StringPtr(data.EncryptionKey)),
		RedirectUris:           redirectURIRequests,
	}

	// Only send grant_types when explicitly set: the API rejects an empty list and otherwise
	// derives the value from the rest of the provider configuration. Note this is the one
	// list here that must *not* go out as a present-but-empty array, so discovery #6's rule
	// is deliberately not applied to it.
	if len(grantTypes) > 0 {
		body.GrantTypes = helpers.CastSliceString[api.GrantTypesEnum](grantTypes)
	}

	return body, diags
}

// fromAPI deliberately never assigns JWKSSources: it is one of the 18 H4 attributes, and the
// inert kind - SDKv2 neither sent it in a request nor read it back, so it only ever holds
// what config put there. Callers seed data from req.Plan or req.State first, so leaving it
// alone preserves that.
func (r *providerOAuth2Resource) fromAPI(ctx context.Context, data *providerOAuth2Model, res *api.OAuth2Provider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.AuthorizationFlow = types.StringValue(res.AuthorizationFlow)
	data.InvalidationFlow = types.StringValue(res.InvalidationFlow)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	jwtFederationSources, d := helpers.MergeStringList(ctx, data.JWTFederationSources, res.JwtFederationSources)
	diags.Append(d...)
	data.JWTFederationSources = jwtFederationSources

	jwtFederationProviders, d := helpers.MergeInt32List(ctx, data.JWTFederationProviders, res.JwtFederationProviders)
	diags.Append(d...)
	data.JWTFederationProviders = jwtFederationProviders

	// grant_types is Optional+Computed, so it is never null in state once applied; the API
	// always returns the derived value.
	grantTypes := make([]string, len(res.GrantTypes))
	for i, gt := range res.GrantTypes {
		grantTypes[i] = string(gt)
	}
	mergedGrantTypes, d := helpers.MergeStringList(ctx, data.GrantTypes, grantTypes)
	diags.Append(d...)
	data.GrantTypes = mergedGrantTypes

	// allowed_redirect_uris keeps SDKv2's merge semantics exactly: entries the API echoes
	// back unchanged hold their configured position, and anything new is appended.
	priorURIs, d := redirectURIsFromList(ctx, data.AllowedRedirectURIs)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	mergedURIs := helpers.ListConsistentMerge(priorURIs, redirectURIsFromAPI(res.RedirectUris))
	if data.AllowedRedirectURIs.IsNull() && len(mergedURIs) == 0 {
		// H1: a never-configured list stays null rather than becoming [].
		data.AllowedRedirectURIs = types.ListNull(redirectURIElementType)
	} else {
		list, d := redirectURIsToList(mergedURIs)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		data.AllowedRedirectURIs = list
	}

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.SigningKey = helpers.StringPtrOrNull(res.SigningKey.Get())
	data.EncryptionKey = helpers.StringPtrOrNull(res.EncryptionKey.Get())
	// No Default, so prior-aware.
	data.LogoutURI = helpers.StringOrNull(data.LogoutURI, res.GetLogoutUri())
	// Required / Computed, so never null.
	data.ClientID = types.StringValue(res.GetClientId())
	data.ClientSecret = types.StringValue(res.GetClientSecret())
	// Have Defaults, so verbatim.
	data.ClientType = types.StringValue(string(res.GetClientType()))
	data.AccessCodeValidity = types.StringValue(res.GetAccessCodeValidity())
	data.AccessTokenValidity = types.StringValue(res.GetAccessTokenValidity())
	data.RefreshTokenValidity = types.StringValue(res.GetRefreshTokenValidity())
	data.RefreshTokenThreshold = types.StringValue(res.GetRefreshTokenThreshold())
	data.IncludeClaimsInIDToken = types.BoolValue(res.GetIncludeClaimsInIdToken())
	data.LogoutMethod = types.StringValue(string(res.GetLogoutMethod()))
	data.SubMode = types.StringValue(string(res.GetSubMode()))
	data.IssuerMode = types.StringValue(string(res.GetIssuerMode()))

	return diags
}

func (r *providerOAuth2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerOAuth2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersOauth2Create(ctx).OAuth2ProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerOAuth2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries jwks_sources - see fromAPI.
	var data providerOAuth2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersOauth2Retrieve(ctx, id).Execute()
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

func (r *providerOAuth2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerOAuth2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersOauth2Update(ctx, id).OAuth2ProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerOAuth2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerOAuth2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersOauth2Destroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
