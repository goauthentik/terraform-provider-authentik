package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	_ resource.Resource                = &providerSCIMResource{}
	_ resource.ResourceWithConfigure   = &providerSCIMResource{}
	_ resource.ResourceWithImportState = &providerSCIMResource{}
)

func newProviderSCIMResource() resource.Resource {
	return &providerSCIMResource{}
}

type providerSCIMResource struct {
	resourceBase
}

type providerSCIMModel struct {
	ID                                types.String         `tfsdk:"id"`
	Name                              types.String         `tfsdk:"name"`
	DryRun                            types.Bool           `tfsdk:"dry_run"`
	URL                               types.String         `tfsdk:"url"`
	Token                             types.String         `tfsdk:"token"`
	AuthMode                          types.String         `tfsdk:"auth_mode"`
	AuthOAuth                         types.String         `tfsdk:"auth_oauth"`
	AuthOAuthParams                   jsontypes.Normalized `tfsdk:"auth_oauth_params"`
	CompatibilityMode                 types.String         `tfsdk:"compatibility_mode"`
	PropertyMappings                  types.List           `tfsdk:"property_mappings"`
	PropertyMappingsGroup             types.List           `tfsdk:"property_mappings_group"`
	ExcludeUsersServiceAccount        types.Bool           `tfsdk:"exclude_users_service_account"`
	GroupFilters                      types.List           `tfsdk:"group_filters"`
	ServiceProviderConfigCacheTimeout types.String         `tfsdk:"service_provider_config_cache_timeout"`
	SyncPageTimeout                   types.String         `tfsdk:"sync_page_timeout"`
	SyncPageSize                      types.Int32          `tfsdk:"sync_page_size"`
}

func (r *providerSCIMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_scim"
}

func (r *providerSCIMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	dryRunDefault := helpers.BoolDefault(false)
	authModeDefault := helpers.StringDefault(string(api.SCIMAUTHENTICATIONMODEENUM_TOKEN))
	authOAuthParamsDefault := helpers.StringDefault("{}")
	compatibilityModeDefault := helpers.StringDefault(string(api.COMPATIBILITYMODEENUM_DEFAULT))
	serviceProviderConfigCacheTimeoutDefault := helpers.StringDefault("hours=1")
	syncPageTimeoutDefault := helpers.StringDefault("minutes=30")
	syncPageSizeDefault := helpers.Int32Default(100)

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
			"dry_run": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             dryRunDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(dryRunDefault.Value())),
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			// token is Sensitive but the API does return it, so this is not an H4 case.
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"auth_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             authModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedSCIMAuthenticationModeEnumEnumValues), helpers.WithDefault(authModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedSCIMAuthenticationModeEnumEnumValues),
				},
			},
			"auth_oauth": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Slug of an OAuth source used for authentication",
			},
			"auth_oauth_params": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				Default:             authOAuthParamsDefault,
				MarkdownDescription: helpers.Desc(helpers.JSONDescription, helpers.WithDefault(authOAuthParamsDefault.Value())),
			},
			"compatibility_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             compatibilityModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedCompatibilityModeEnumEnumValues), helpers.WithDefault(compatibilityModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedCompatibilityModeEnumEnumValues),
				},
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"property_mappings_group": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"exclude_users_service_account": schema.BoolAttribute{
				Optional: true,
			},
			"group_filters": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"service_provider_config_cache_timeout": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             serviceProviderConfigCacheTimeoutDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(serviceProviderConfigCacheTimeoutDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"sync_page_timeout": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncPageTimeoutDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(syncPageTimeoutDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"sync_page_size": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             syncPageSizeDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(syncPageSizeDefault.Value())),
			},
		},
	}
}

func (r *providerSCIMResource) toRequest(ctx context.Context, data *providerSCIMModel) (*api.SCIMProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	groupFilters, d := helpers.SliceOrEmpty[string](ctx, data.GroupFilters)
	diags.Append(d...)

	var authOAuthParams map[string]any
	diags.Append(data.AuthOAuthParams.Unmarshal(&authOAuthParams)...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.SCIMProviderRequest{
		Name:                              data.Name.ValueString(),
		Url:                               data.URL.ValueString(),
		AuthMode:                          api.SCIMAuthenticationModeEnum(data.AuthMode.ValueString()).Ptr(),
		AuthOauth:                         *api.NewNullableString(helpers.StringPtr(data.AuthOAuth)),
		Token:                             helpers.StringPtr(data.Token),
		PropertyMappings:                  propertyMappings,
		PropertyMappingsGroup:             propertyMappingsGroup,
		ExcludeUsersServiceAccount:        new(data.ExcludeUsersServiceAccount.ValueBool()),
		CompatibilityMode:                 api.CompatibilityModeEnum(data.CompatibilityMode.ValueString()).Ptr(),
		GroupFilters:                      groupFilters,
		DryRun:                            new(data.DryRun.ValueBool()),
		ServiceProviderConfigCacheTimeout: helpers.StringPtr(data.ServiceProviderConfigCacheTimeout),
		SyncPageTimeout:                   helpers.StringPtr(data.SyncPageTimeout),
		SyncPageSize:                      helpers.Int32Ptr(data.SyncPageSize),
		AuthOauthParams:                   authOAuthParams,
	}, diags
}

func (r *providerSCIMResource) fromAPI(ctx context.Context, data *providerSCIMModel, res *api.SCIMProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.URL = types.StringValue(res.Url)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	propertyMappingsGroup, d := helpers.MergeStringList(ctx, data.PropertyMappingsGroup, res.PropertyMappingsGroup)
	diags.Append(d...)
	data.PropertyMappingsGroup = propertyMappingsGroup

	groupFilters, d := helpers.MergeStringList(ctx, data.GroupFilters, res.GroupFilters)
	diags.Append(d...)
	data.GroupFilters = groupFilters

	// Nullable API field.
	data.AuthOAuth = helpers.StringPtrOrNull(res.AuthOauth.Get())
	// No Defaults, so prior-aware.
	data.Token = helpers.StringOrNull(data.Token, res.GetToken())
	data.ExcludeUsersServiceAccount = helpers.BoolOrNull(data.ExcludeUsersServiceAccount, res.GetExcludeUsersServiceAccount())
	// Have Defaults, so verbatim.
	data.DryRun = types.BoolValue(res.GetDryRun())
	data.AuthMode = types.StringValue(string(res.GetAuthMode()))
	data.CompatibilityMode = types.StringValue(string(res.GetCompatibilityMode()))
	data.ServiceProviderConfigCacheTimeout = types.StringValue(res.GetServiceProviderConfigCacheTimeout())
	data.SyncPageTimeout = types.StringValue(res.GetSyncPageTimeout())
	data.SyncPageSize = types.Int32Value(res.GetSyncPageSize())

	paramsBytes, err := json.Marshal(res.AuthOauthParams)
	if err != nil {
		diags.AddError("Failed to encode SCIM provider auth_oauth_params", err.Error())
		return diags
	}
	data.AuthOAuthParams = jsontypes.NewNormalizedValue(string(paramsBytes))

	return diags
}

func (r *providerSCIMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerSCIMModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersScimCreate(ctx).SCIMProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerSCIMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerSCIMModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersScimRetrieve(ctx, id).Execute()
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

func (r *providerSCIMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerSCIMModel
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

	res, hr, err := r.client.ProvidersAPI.ProvidersScimUpdate(ctx, id).SCIMProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerSCIMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerSCIMModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersScimDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
