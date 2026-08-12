package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ datasource.DataSource                     = &groupDataSource{}
	_ datasource.DataSourceWithConfigure        = &groupDataSource{}
	_ datasource.DataSourceWithConfigValidators = &groupDataSource{}
)

func newGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	dataSourceBase
}

// groupMemberAttrTypes is shared between the schema (as an ObjectType's AttrTypes) and
// types.ListValueFrom (to build the users_obj list), so the two can't drift apart.
var groupMemberAttrTypes = map[string]attr.Type{
	"pk":         types.Int32Type,
	"username":   types.StringType,
	"name":       types.StringType,
	"is_active":  types.BoolType,
	"last_login": types.StringType,
	"email":      types.StringType,
	"attributes": types.StringType,
	"uid":        types.StringType,
}

type groupMemberModel struct {
	Pk         types.Int32  `tfsdk:"pk"`
	Username   types.String `tfsdk:"username"`
	Name       types.String `tfsdk:"name"`
	IsActive   types.Bool   `tfsdk:"is_active"`
	LastLogin  types.String `tfsdk:"last_login"`
	Email      types.String `tfsdk:"email"`
	Attributes types.String `tfsdk:"attributes"`
	UID        types.String `tfsdk:"uid"`
}

type groupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Pk           types.String `tfsdk:"pk"`
	NumPk        types.Int32  `tfsdk:"num_pk"`
	Name         types.String `tfsdk:"name"`
	IncludeUsers types.Bool   `tfsdk:"include_users"`
	IsSuperuser  types.Bool   `tfsdk:"is_superuser"`
	Parents      types.List   `tfsdk:"parents"`
	ParentName   types.String `tfsdk:"parent_name"`
	Users        types.List   `tfsdk:"users"`
	Attributes   types.String `tfsdk:"attributes"`
	UsersObj     types.List   `tfsdk:"users_obj"`
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema is byte-for-byte equivalent to pkg/sdkprovider's dataSourceGroup, plus the
// explicit "id" attribute the framework doesn't inject automatically (H3 for data
// sources - only 6 of 25 SDKv2 data sources declare it today, this one doesn't).
// users_obj uses ListAttribute+ObjectType rather than ListNestedAttribute: SDKv2's
// CoreConfigSchema routes Computed-only nested lists (Computed && !Optional, which is
// every nested list in a data source) to protocol attributes, not blocks, so
// ListNestedAttribute here would emit a different (NestedType) protocol schema.
func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Directory --- Get groups by pk or name",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"pk": schema.StringAttribute{
				Optional: true,
			},
			"num_pk": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"include_users": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: helpers.Desc(
					"Whether to include group members. Note that depending on group size, this can make the Terraform state a lot larger.",
					helpers.WithDefault(true),
				),
			},
			"is_superuser": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"parents": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"parent_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"users": schema.ListAttribute{
				ElementType:         types.Int32Type,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"attributes": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
			"users_obj": schema.ListAttribute{
				ElementType:         types.ObjectType{AttrTypes: groupMemberAttrTypes},
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
			},
		},
	}
}

func (d *groupDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("pk"),
			path.MatchRoot("name"),
		),
	}
}

func memberFromAPI(member api.PartialUser) (groupMemberModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	lastLogin := ""
	if t, ok := member.GetLastLoginOk(); ok && t != nil {
		if b, err := t.MarshalText(); err == nil && b != nil {
			lastLogin = string(b)
		}
	}

	attrBytes, err := json.Marshal(member.GetAttributes())
	if err != nil {
		diags.AddError("Failed to encode group member attributes", err.Error())
		return groupMemberModel{}, diags
	}

	return groupMemberModel{
		Pk:         types.Int32Value(member.GetPk()),
		Username:   types.StringValue(member.GetUsername()),
		Name:       types.StringValue(member.GetName()),
		IsActive:   types.BoolValue(member.GetIsActive()),
		LastLogin:  types.StringValue(lastLogin),
		Email:      types.StringValue(member.GetEmail()),
		Attributes: types.StringValue(string(attrBytes)),
		UID:        types.StringValue(member.GetUid()),
	}, diags
}

func groupDataSourceFromAPI(ctx context.Context, group *api.Group) (groupDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrBytes, err := json.Marshal(group.GetAttributes())
	if err != nil {
		diags.AddError("Failed to encode group attributes", err.Error())
		return groupDataSourceModel{}, diags
	}

	parents, d := types.ListValueFrom(ctx, types.StringType, group.GetParents())
	diags.Append(d...)

	users, d := types.ListValueFrom(ctx, types.Int32Type, group.GetUsers())
	diags.Append(d...)

	members := make([]groupMemberModel, 0, len(group.GetUsersObj()))
	for _, m := range group.GetUsersObj() {
		member, d := memberFromAPI(m)
		diags.Append(d...)
		members = append(members, member)
	}
	usersObj, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: groupMemberAttrTypes}, members)
	diags.Append(d...)

	if diags.HasError() {
		return groupDataSourceModel{}, diags
	}

	return groupDataSourceModel{
		ID:          types.StringValue(group.GetPk()),
		Pk:          types.StringValue(group.GetPk()),
		NumPk:       types.Int32Value(group.GetNumPk()),
		Name:        types.StringValue(group.GetName()),
		IsSuperuser: types.BoolValue(group.GetIsSuperuser()),
		Parents:     parents,
		// parent_name is declared but never populated in the SDKv2 data source either -
		// carried over as-is rather than fixed as part of this migration.
		ParentName: types.StringValue(""),
		Users:      users,
		Attributes: types.StringValue(string(attrBytes)),
		UsersObj:   usersObj,
	}, diags
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer d.span(ctx, "read")()

	var data groupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// datasource/schema has no Default mechanism (unlike resource/schema), so
	// include_users' SDKv2 Default: true is applied here instead.
	includeUsers := true
	if !data.IncludeUsers.IsNull() {
		includeUsers = data.IncludeUsers.ValueBool()
	}

	var res *api.Group
	switch {
	case !data.Pk.IsNull():
		r, hr, err := d.client.CoreAPI.CoreGroupsRetrieve(ctx, data.Pk.ValueString()).IncludeUsers(includeUsers).Execute()
		if err != nil {
			resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
			return
		}
		res = r
	case !data.Name.IsNull():
		list, hr, err := d.client.CoreAPI.CoreGroupsList(ctx).IncludeUsers(includeUsers).Name(data.Name.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
			return
		}
		if len(list.Results) < 1 {
			resp.Diagnostics.AddError("No matching groups found", "")
			return
		}
		if len(list.Results) > 1 {
			resp.Diagnostics.AddError("Multiple groups found", "")
			return
		}
		res = &list.Results[0]
	default:
		resp.Diagnostics.AddError("Neither pk nor name were provided", "")
		return
	}

	model, diags := groupDataSourceFromAPI(ctx, res)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	model.IncludeUsers = types.BoolValue(includeUsers)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
