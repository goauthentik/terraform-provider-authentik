package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func newGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	resourceBase
}

type groupModel struct {
	ID          types.String         `tfsdk:"id"`
	Name        types.String         `tfsdk:"name"`
	IsSuperuser types.Bool           `tfsdk:"is_superuser"`
	Parents     types.List           `tfsdk:"parents"`
	Users       types.List           `tfsdk:"users"`
	Roles       types.List           `tfsdk:"roles"`
	Attributes  jsontypes.Normalized `tfsdk:"attributes"`
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema is byte-for-byte equivalent to pkg/sdkprovider's resourceGroup, plus the
// explicit "id" attribute the framework doesn't inject automatically (H3). Pk is a
// stable, server-generated key (never derived from a mutable attribute), so
// UseStateForUnknown is safe here - contrast H3's ten slug/identifier-keyed resources,
// which must not have it.
func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributesDefault := helpers.StringDefault("{}")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Directory --- ",
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
			"is_superuser": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             helpers.BoolDefault(false),
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(false)),
			},
			"parents": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"users": schema.ListAttribute{
				ElementType:         types.Int32Type,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"attributes": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				Default:             attributesDefault,
				MarkdownDescription: helpers.Desc(helpers.JSONDescription, helpers.WithDefault(attributesDefault.Value())),
			},
			"roles": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

// toRequest converts the plan/state model to an API request body. Unlike SDKv2's
// GetJSON/CastSlice helpers, ElementsAs on a null types.List sets the target to nil
// without error, so no null-guard is needed before it.
func (r *groupResource) toRequest(ctx context.Context, data *groupModel) (*api.GroupRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var parents []string
	diags.Append(data.Parents.ElementsAs(ctx, &parents, false)...)

	var users []int32
	diags.Append(data.Users.ElementsAs(ctx, &users, false)...)

	var roles []string
	diags.Append(data.Roles.ElementsAs(ctx, &roles, false)...)

	var attributes map[string]any
	diags.Append(data.Attributes.Unmarshal(&attributes)...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.GroupRequest{
		Name:        data.Name.ValueString(),
		IsSuperuser: new(data.IsSuperuser.ValueBool()),
		Parents:     parents,
		Users:       users,
		Roles:       roles,
		Attributes:  attributes,
	}, diags
}

// fromAPI is the H4/H1 fix: data is seeded from the plan (Create/Update) or prior state
// (Read) before this overwrites it, so list ordering is preserved via
// helpers.MergeStringList/MergeInt32List (the executable spec for that is
// TestResourceGroupReadRolesPreserveConfiguredOrder in resource_group_test.go) and
// nothing marshals a zero value where the config said null.
func (r *groupResource) fromAPI(ctx context.Context, data *groupModel, res *api.Group) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.IsSuperuser = helpers.BoolOrNull(data.IsSuperuser, res.GetIsSuperuser())

	parents, d := helpers.MergeStringList(ctx, data.Parents, res.Parents)
	diags.Append(d...)
	data.Parents = parents

	users, d := helpers.MergeInt32List(ctx, data.Users, res.Users)
	diags.Append(d...)
	data.Users = users

	roles, d := helpers.MergeStringList(ctx, data.Roles, res.Roles)
	diags.Append(d...)
	data.Roles = roles

	attrBytes, err := json.Marshal(res.Attributes)
	if err != nil {
		diags.AddError("Failed to encode group attributes", err.Error())
		return diags
	}
	data.Attributes = jsontypes.NewNormalizedValue(string(attrBytes))

	return diags
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data groupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.CoreAPI.CoreGroupsCreate(ctx).GroupRequest(*body).Execute()
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

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.CoreAPI.CoreGroupsRetrieve(ctx, data.ID.ValueString()).IncludeUsers(false).Execute()
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

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data groupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// data.ID is already the existing group's Pk: "id" has UseStateForUnknown, so the
	// plan value is the prior state value carried forward, not unknown.
	res, hr, err := r.client.CoreAPI.CoreGroupsUpdate(ctx, data.ID.ValueString()).GroupRequest(*body).Execute()
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

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.CoreAPI.CoreGroupsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
