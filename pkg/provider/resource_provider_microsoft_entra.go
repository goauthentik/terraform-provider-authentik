package provider

import (
	"context"
	"strconv"

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
	_ resource.Resource                = &providerMicrosoftEntraResource{}
	_ resource.ResourceWithConfigure   = &providerMicrosoftEntraResource{}
	_ resource.ResourceWithImportState = &providerMicrosoftEntraResource{}
)

func newProviderMicrosoftEntraResource() resource.Resource {
	return &providerMicrosoftEntraResource{}
}

type providerMicrosoftEntraResource struct {
	resourceBase
}

type providerMicrosoftEntraModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	DryRun                     types.Bool   `tfsdk:"dry_run"`
	ClientID                   types.String `tfsdk:"client_id"`
	ClientSecret               types.String `tfsdk:"client_secret"`
	TenantID                   types.String `tfsdk:"tenant_id"`
	PropertyMappings           types.List   `tfsdk:"property_mappings"`
	PropertyMappingsGroup      types.List   `tfsdk:"property_mappings_group"`
	ExcludeUsersServiceAccount types.Bool   `tfsdk:"exclude_users_service_account"`
	FilterGroup                types.String `tfsdk:"filter_group"`
	UserDeleteAction           types.String `tfsdk:"user_delete_action"`
	GroupDeleteAction          types.String `tfsdk:"group_delete_action"`
	SyncPageTimeout            types.String `tfsdk:"sync_page_timeout"`
	SyncPageSize               types.Int32  `tfsdk:"sync_page_size"`
}

func (r *providerMicrosoftEntraResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_microsoft_entra"
}

func (r *providerMicrosoftEntraResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	dryRunDefault := helpers.BoolDefault(false)
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
			"client_id": schema.StringAttribute{
				Required: true,
			},
			// client_secret is Required and Sensitive but the API does return it, so this is
			// not an H4 case and fromAPI writes it normally.
			"client_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"tenant_id": schema.StringAttribute{
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
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(outgoingSyncGroupDeleteActions), helpers.WithDefault(userDeleteActionDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(outgoingSyncGroupDeleteActions),
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

func (r *providerMicrosoftEntraResource) toRequest(ctx context.Context, data *providerMicrosoftEntraModel) (*api.MicrosoftEntraProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.MicrosoftEntraProviderRequest{
		Name:                       data.Name.ValueString(),
		ClientId:                   data.ClientID.ValueString(),
		ClientSecret:               data.ClientSecret.ValueString(),
		TenantId:                   data.TenantID.ValueString(),
		PropertyMappings:           propertyMappings,
		PropertyMappingsGroup:      propertyMappingsGroup,
		ExcludeUsersServiceAccount: new(data.ExcludeUsersServiceAccount.ValueBool()),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(data.UserDeleteAction.ValueString()).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(data.GroupDeleteAction.ValueString()).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.StringPtr(data.FilterGroup)),
		DryRun:                     new(data.DryRun.ValueBool()),
		SyncPageTimeout:            helpers.StringPtr(data.SyncPageTimeout),
		SyncPageSize:               helpers.Int32Ptr(data.SyncPageSize),
	}, diags
}

func (r *providerMicrosoftEntraResource) fromAPI(ctx context.Context, data *providerMicrosoftEntraModel, res *api.MicrosoftEntraProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.ClientID = types.StringValue(res.ClientId)
	data.ClientSecret = types.StringValue(res.ClientSecret)
	data.TenantID = types.StringValue(res.TenantId)

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

	return diags
}

func (r *providerMicrosoftEntraResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerMicrosoftEntraModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersMicrosoftEntraCreate(ctx).MicrosoftEntraProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerMicrosoftEntraResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerMicrosoftEntraModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersMicrosoftEntraRetrieve(ctx, id).Execute()
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

func (r *providerMicrosoftEntraResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerMicrosoftEntraModel
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

	res, hr, err := r.client.ProvidersAPI.ProvidersMicrosoftEntraUpdate(ctx, id).MicrosoftEntraProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerMicrosoftEntraResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerMicrosoftEntraModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersMicrosoftEntraDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
