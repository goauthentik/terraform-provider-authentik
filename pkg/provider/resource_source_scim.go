package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &sourceSCIMResource{}
	_ resource.ResourceWithConfigure   = &sourceSCIMResource{}
	_ resource.ResourceWithImportState = &sourceSCIMResource{}
)

func newSourceSCIMResource() resource.Resource {
	return &sourceSCIMResource{}
}

type sourceSCIMResource struct {
	resourceBase
}

type sourceSCIMModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	UUID                  types.String `tfsdk:"uuid"`
	Slug                  types.String `tfsdk:"slug"`
	UserPathTemplate      types.String `tfsdk:"user_path_template"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	PropertyMappings      types.List   `tfsdk:"property_mappings"`
	PropertyMappingsGroup types.List   `tfsdk:"property_mappings_group"`
	SCIMURL               types.String `tfsdk:"scim_url"`
	Token                 types.String `tfsdk:"token"`
}

func (r *sourceSCIMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_scim"
}

func (r *sourceSCIMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Directory --- ",
		Attributes: map[string]schema.Attribute{
			// H3: this resource's id is res.Slug, which is derived from the mutable slug
			// attribute, so it must NOT carry UseStateForUnknown - pinning the old value
			// there would make apply fail whenever the slug changes. The cost is that id
			// shows "(known after apply)" on every update, which is the accepted tradeoff
			// documented for all ten H3 resources.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			// uuid is res.Pk: server-generated once at creation and stable thereafter, so
			// unlike id it is safe to pin with UseStateForUnknown.
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
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enabledDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enabledDefault.Value())),
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"property_mappings_group": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			// Both of these are Computed-only and derived from the source's slug, so per the
			// UseStateForUnknown rule ("safe iff the server value is a pure function of this
			// resource's own config plus immutable creation-time state") they must not pin
			// prior state either - a slug change changes root_url.
			//
			// The "SCIM URL" description on token is a copy-paste slip in the SDKv2 schema,
			// but it is in the published docs, so changing it here would fail the
			// docs-drift gate. Left alone deliberately; worth fixing as its own change.
			"scim_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("SCIM URL", helpers.Generated()),
			},
			"token": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("SCIM URL", helpers.Generated()),
			},
		},
	}
}

func (r *sourceSCIMResource) toRequest(ctx context.Context, data *sourceSCIMModel) (*api.SCIMSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.SCIMSourceRequest{
		Name:                  data.Name.ValueString(),
		Slug:                  data.Slug.ValueString(),
		Enabled:               new(data.Enabled.ValueBool()),
		UserPathTemplate:      new(data.UserPathTemplate.ValueString()),
		UserPropertyMappings:  propertyMappings,
		GroupPropertyMappings: propertyMappingsGroup,
	}, diags
}

func (r *sourceSCIMResource) fromAPI(ctx context.Context, data *sourceSCIMModel, res *api.SCIMSource) diag.Diagnostics {
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

	// Both have Defaults, so they take the API value verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())

	// SDKv2 fetched the source a second time here, with an identical
	// SourcesScimRetrieve call, purely to read these two fields off the response it
	// already had. Collapsed into the single call above.
	data.SCIMURL = types.StringValue(res.RootUrl)
	data.Token = types.StringValue(res.TokenObj.Identifier)

	return diags
}

func (r *sourceSCIMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceSCIMModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesScimCreate(ctx).SCIMSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceSCIMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data sourceSCIMModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesScimRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *sourceSCIMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceSCIMModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// H3/discovery #4: the plan's id is the unknown placeholder whenever the slug changes,
	// so the URL has to use the prior id from state.
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

	res, hr, err := r.client.SourcesAPI.SourcesScimUpdate(ctx, priorID.ValueString()).SCIMSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceSCIMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceSCIMModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesScimDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
