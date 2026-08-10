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
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithConfigure   = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

func newApplicationResource() resource.Resource {
	return &applicationResource{}
}

type applicationResource struct {
	resourceBase
}

type applicationModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Group                types.String `tfsdk:"group"`
	UUID                 types.String `tfsdk:"uuid"`
	Slug                 types.String `tfsdk:"slug"`
	ProtocolProvider     types.Int32  `tfsdk:"protocol_provider"`
	BackchannelProviders types.List   `tfsdk:"backchannel_providers"`
	MetaLaunchURL        types.String `tfsdk:"meta_launch_url"`
	MetaIcon             types.String `tfsdk:"meta_icon"`
	MetaDescription      types.String `tfsdk:"meta_description"`
	MetaPublisher        types.String `tfsdk:"meta_publisher"`
	PolicyEngineMode     types.String `tfsdk:"policy_engine_mode"`
	OpenInNewTab         types.Bool   `tfsdk:"open_in_new_tab"`
	MetaHide             types.Bool   `tfsdk:"meta_hide"`
}

func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema is byte-for-byte equivalent to pkg/sdkprovider's resourceApplication. "id" has
// no UseStateForUnknown (H3): unlike authentik_group's stable Pk, an application's id is
// its slug, which is Required config and can change - pinning the plan value there would
// make apply fail whenever slug changes. "uuid" is the opposite case: server-generated
// and stable once created, so it gets UseStateForUnknown same as authentik_group's id.
func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	openInNewTabDefault := helpers.BoolDefault(false)
	metaHideDefault := helpers.BoolDefault(false)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Applications --- ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"group": schema.StringAttribute{
				Optional: true,
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
			"protocol_provider": schema.Int32Attribute{
				Optional: true,
			},
			"backchannel_providers": schema.ListAttribute{
				ElementType: types.Int32Type,
				Optional:    true,
			},
			"meta_launch_url": schema.StringAttribute{
				Optional: true,
			},
			"meta_icon": schema.StringAttribute{
				Optional: true,
			},
			"meta_description": schema.StringAttribute{
				Optional: true,
			},
			"meta_publisher": schema.StringAttribute{
				Optional: true,
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
			"open_in_new_tab": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             openInNewTabDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(openInNewTabDefault.Value())),
			},
			"meta_hide": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             metaHideDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(metaHideDefault.Value())),
			},
		},
	}
}

func (r *applicationResource) toRequest(ctx context.Context, data *applicationModel) (*api.ApplicationRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The non-nil guard this used to spell out inline now lives in helpers.SliceOrEmpty,
	// which also handles the unknown case.
	backchannelProviders, d := helpers.SliceOrEmpty[int32](ctx, data.BackchannelProviders)
	diags.Append(d...)

	return &api.ApplicationRequest{
		Name:                 data.Name.ValueString(),
		Slug:                 data.Slug.ValueString(),
		Provider:             *api.NewNullableInt32(helpers.Int32Ptr(data.ProtocolProvider)),
		BackchannelProviders: backchannelProviders,
		OpenInNewTab:         new(data.OpenInNewTab.ValueBool()),
		MetaLaunchUrl:        helpers.StringPtrEmpty(data.MetaLaunchURL),
		MetaIcon:             helpers.StringPtrEmpty(data.MetaIcon),
		MetaDescription:      helpers.StringPtr(data.MetaDescription),
		MetaPublisher:        helpers.StringPtr(data.MetaPublisher),
		PolicyEngineMode:     api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		Group:                helpers.StringPtrEmpty(data.Group),
		MetaHide:             new(data.MetaHide.ValueBool()),
	}, diags
}

func (r *applicationResource) fromAPI(ctx context.Context, data *applicationModel, res *api.Application) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Group = helpers.StringOrNull(data.Group, res.GetGroup())
	data.Slug = types.StringValue(res.Slug)
	data.OpenInNewTab = types.BoolValue(res.GetOpenInNewTab())
	data.MetaHide = types.BoolValue(res.GetMetaHide())
	data.ProtocolProvider = helpers.Int32PtrOrNull(res.Provider.Get())
	data.MetaLaunchURL = helpers.StringOrNull(data.MetaLaunchURL, res.GetMetaLaunchUrl())
	data.MetaIcon = helpers.StringOrNull(data.MetaIcon, res.GetMetaIcon())
	data.MetaDescription = helpers.StringOrNull(data.MetaDescription, res.GetMetaDescription())
	data.MetaPublisher = helpers.StringOrNull(data.MetaPublisher, res.GetMetaPublisher())
	if res.PolicyEngineMode != nil {
		data.PolicyEngineMode = types.StringValue(string(*res.PolicyEngineMode))
	}

	backchannelProviders, d := helpers.MergeInt32List(ctx, data.BackchannelProviders, res.BackchannelProviders)
	diags.Append(d...)
	data.BackchannelProviders = backchannelProviders

	return diags
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data applicationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.CoreAPI.CoreApplicationsCreate(ctx).ApplicationRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.CoreAPI.CoreApplicationsRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data applicationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	res, hr, err := r.client.CoreAPI.CoreApplicationsUpdate(ctx, priorID.ValueString()).ApplicationRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.CoreAPI.CoreApplicationsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
