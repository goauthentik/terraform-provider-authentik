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
	_ resource.Resource                = &providerGoogleWorkspaceResource{}
	_ resource.ResourceWithConfigure   = &providerGoogleWorkspaceResource{}
	_ resource.ResourceWithImportState = &providerGoogleWorkspaceResource{}
)

// outgoingSyncGroupDeleteActions is the subset of api.OutgoingSyncDeleteAction the
// *group* delete action accepts (groups cannot be suspended). SDKv2 spelled the list out
// twice per resource, once for the description and once for the validator.
// The order matters: EnumToDescription renders the list verbatim into the docs, so it must
// match the order SDKv2 spelled out (delete first, then do_nothing).
var outgoingSyncGroupDeleteActions = []api.OutgoingSyncDeleteAction{
	api.OUTGOINGSYNCDELETEACTION_DELETE,
	api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
}

func newProviderGoogleWorkspaceResource() resource.Resource {
	return &providerGoogleWorkspaceResource{}
}

type providerGoogleWorkspaceResource struct {
	resourceBase
}

type providerGoogleWorkspaceModel struct {
	ID                         types.String         `tfsdk:"id"`
	Name                       types.String         `tfsdk:"name"`
	DryRun                     types.Bool           `tfsdk:"dry_run"`
	Credentials                jsontypes.Normalized `tfsdk:"credentials"`
	DelegatedSubject           types.String         `tfsdk:"delegated_subject"`
	DefaultGroupEmailDomain    types.String         `tfsdk:"default_group_email_domain"`
	PropertyMappings           types.List           `tfsdk:"property_mappings"`
	PropertyMappingsGroup      types.List           `tfsdk:"property_mappings_group"`
	ExcludeUsersServiceAccount types.Bool           `tfsdk:"exclude_users_service_account"`
	FilterGroup                types.String         `tfsdk:"filter_group"`
	UserDeleteAction           types.String         `tfsdk:"user_delete_action"`
	GroupDeleteAction          types.String         `tfsdk:"group_delete_action"`
	SyncPageTimeout            types.String         `tfsdk:"sync_page_timeout"`
	SyncPageSize               types.Int32          `tfsdk:"sync_page_size"`
}

func (r *providerGoogleWorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_google_workspace"
}

func (r *providerGoogleWorkspaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	dryRunDefault := helpers.BoolDefault(false)
	credentialsDefault := helpers.StringDefault("{}")
	userDeleteActionDefault := helpers.StringDefault(string(api.OUTGOINGSYNCDELETEACTION_DELETE))
	groupDeleteActionDefault := helpers.StringDefault(string(api.OUTGOINGSYNCDELETEACTION_DELETE))
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
			"credentials": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				Default:             credentialsDefault,
				MarkdownDescription: helpers.Desc(helpers.JSONDescription, helpers.WithDefault(credentialsDefault.Value())),
			},
			"delegated_subject": schema.StringAttribute{
				Optional: true,
			},
			"default_group_email_domain": schema.StringAttribute{
				Required: true,
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
			"filter_group": schema.StringAttribute{
				Optional: true,
			},
			"user_delete_action": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userDeleteActionDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedOutgoingSyncDeleteActionEnumValues), helpers.WithDefault(userDeleteActionDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedOutgoingSyncDeleteActionEnumValues),
				},
			},
			"group_delete_action": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             groupDeleteActionDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(outgoingSyncGroupDeleteActions), helpers.WithDefault(groupDeleteActionDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(outgoingSyncGroupDeleteActions),
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

func (r *providerGoogleWorkspaceResource) toRequest(ctx context.Context, data *providerGoogleWorkspaceModel) (*api.GoogleWorkspaceProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	var credentials map[string]any
	diags.Append(data.Credentials.Unmarshal(&credentials)...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.GoogleWorkspaceProviderRequest{
		Name:                       data.Name.ValueString(),
		DelegatedSubject:           data.DelegatedSubject.ValueString(),
		DefaultGroupEmailDomain:    data.DefaultGroupEmailDomain.ValueString(),
		PropertyMappings:           propertyMappings,
		PropertyMappingsGroup:      propertyMappingsGroup,
		ExcludeUsersServiceAccount: new(data.ExcludeUsersServiceAccount.ValueBool()),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(data.UserDeleteAction.ValueString()).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(data.GroupDeleteAction.ValueString()).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.StringPtr(data.FilterGroup)),
		DryRun:                     new(data.DryRun.ValueBool()),
		SyncPageTimeout:            helpers.StringPtr(data.SyncPageTimeout),
		SyncPageSize:               helpers.Int32Ptr(data.SyncPageSize),
		Credentials:                credentials,
	}, diags
}

func (r *providerGoogleWorkspaceResource) fromAPI(ctx context.Context, data *providerGoogleWorkspaceModel, res *api.GoogleWorkspaceProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.DelegatedSubject = types.StringValue(res.DelegatedSubject)
	data.DefaultGroupEmailDomain = types.StringValue(res.DefaultGroupEmailDomain)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	propertyMappingsGroup, d := helpers.MergeStringList(ctx, data.PropertyMappingsGroup, res.PropertyMappingsGroup)
	diags.Append(d...)
	data.PropertyMappingsGroup = propertyMappingsGroup

	// filter_group is a NullableString, and SDKv2 passed the wrapper struct itself to
	// SetWrapper rather than calling .Get() on it. d.Set rejects that with
	// "expected type 'string', got unconvertible type 'api.NullableString'", and SetWrapper
	// panics on any error - so every Read of this resource crashed the provider. There is no
	// acceptance test for it, which is why it went unnoticed. Fixed by unwrapping properly.
	data.FilterGroup = helpers.StringPtrOrNull(res.FilterGroup.Get())
	// No Default, so prior-aware.
	data.ExcludeUsersServiceAccount = helpers.BoolOrNull(data.ExcludeUsersServiceAccount, res.GetExcludeUsersServiceAccount())
	// Have Defaults, so verbatim.
	data.DryRun = types.BoolValue(res.GetDryRun())
	data.UserDeleteAction = types.StringValue(string(res.GetUserDeleteAction()))
	data.GroupDeleteAction = types.StringValue(string(res.GetGroupDeleteAction()))
	data.SyncPageTimeout = types.StringValue(res.GetSyncPageTimeout())
	data.SyncPageSize = types.Int32Value(res.GetSyncPageSize())

	credentialsBytes, err := json.Marshal(res.Credentials)
	if err != nil {
		diags.AddError("Failed to encode Google Workspace provider credentials", err.Error())
		return diags
	}
	data.Credentials = jsontypes.NewNormalizedValue(string(credentialsBytes))

	return diags
}

func (r *providerGoogleWorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerGoogleWorkspaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersGoogleWorkspaceCreate(ctx).GoogleWorkspaceProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerGoogleWorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerGoogleWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersGoogleWorkspaceRetrieve(ctx, id).Execute()
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

func (r *providerGoogleWorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerGoogleWorkspaceModel
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

	res, hr, err := r.client.ProvidersAPI.ProvidersGoogleWorkspaceUpdate(ctx, id).GoogleWorkspaceProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerGoogleWorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerGoogleWorkspaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersGoogleWorkspaceDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
