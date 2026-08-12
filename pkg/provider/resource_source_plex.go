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
	_ resource.Resource                = &sourcePlexResource{}
	_ resource.ResourceWithConfigure   = &sourcePlexResource{}
	_ resource.ResourceWithImportState = &sourcePlexResource{}
)

func newSourcePlexResource() resource.Resource {
	return &sourcePlexResource{}
}

type sourcePlexResource struct {
	resourceBase
}

type sourcePlexModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	UUID               types.String `tfsdk:"uuid"`
	Slug               types.String `tfsdk:"slug"`
	UserPathTemplate   types.String `tfsdk:"user_path_template"`
	AuthenticationFlow types.String `tfsdk:"authentication_flow"`
	EnrollmentFlow     types.String `tfsdk:"enrollment_flow"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	Promoted           types.Bool   `tfsdk:"promoted"`
	PolicyEngineMode   types.String `tfsdk:"policy_engine_mode"`
	UserMatchingMode   types.String `tfsdk:"user_matching_mode"`
	GroupMatchingMode  types.String `tfsdk:"group_matching_mode"`
	ClientID           types.String `tfsdk:"client_id"`
	AllowedServers     types.List   `tfsdk:"allowed_servers"`
	AllowFriends       types.Bool   `tfsdk:"allow_friends"`
	PlexToken          types.String `tfsdk:"plex_token"`
}

func (r *sourcePlexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_plex"
}

func (r *sourcePlexResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	promotedDefault := helpers.BoolDefault(false)
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	userMatchingModeDefault := helpers.StringDefault(string(api.USERMATCHINGMODEENUM_IDENTIFIER))
	groupMatchingModeDefault := helpers.StringDefault(string(api.GROUPMATCHINGMODEENUM_IDENTIFIER))
	allowFriendsDefault := helpers.BoolDefault(true)

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
			"client_id": schema.StringAttribute{
				Required: true,
			},
			"allowed_servers": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"allow_friends": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             allowFriendsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(allowFriendsDefault.Value())),
			},
			// Unlike the other sources' secrets, plex_token *is* returned by the API, so
			// this resource is not one of the 18 H4 cases and fromAPI writes it normally.
			"plex_token": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *sourcePlexResource) toRequest(ctx context.Context, data *sourcePlexModel) (*api.PlexSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	allowedServers, d := helpers.SliceOrEmpty[string](ctx, data.AllowedServers)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.PlexSourceRequest{
		Name:               data.Name.ValueString(),
		Slug:               data.Slug.ValueString(),
		Enabled:            new(data.Enabled.ValueBool()),
		Promoted:           new(data.Promoted.ValueBool()),
		UserPathTemplate:   new(data.UserPathTemplate.ValueString()),
		PolicyEngineMode:   api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		UserMatchingMode:   api.UserMatchingModeEnum(data.UserMatchingMode.ValueString()).Ptr(),
		GroupMatchingMode:  api.GroupMatchingModeEnum(data.GroupMatchingMode.ValueString()).Ptr(),
		AuthenticationFlow: *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		EnrollmentFlow:     *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		ClientId:           new(data.ClientID.ValueString()),
		AllowFriends:       new(data.AllowFriends.ValueBool()),
		PlexToken:          data.PlexToken.ValueString(),
		AllowedServers:     allowedServers,
	}, diags
}

func (r *sourcePlexResource) fromAPI(ctx context.Context, data *sourcePlexModel, res *api.PlexSource) diag.Diagnostics {
	var diags diag.Diagnostics

	// H3: id is the slug, not the Pk.
	data.ID = types.StringValue(res.Slug)
	data.Name = types.StringValue(res.Name)
	data.Slug = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)

	allowedServers, d := helpers.MergeStringList(ctx, data.AllowedServers, res.AllowedServers)
	diags.Append(d...)
	data.AllowedServers = allowedServers

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.EnrollmentFlow = helpers.StringPtrOrNull(res.EnrollmentFlow.Get())
	// Required, so never null.
	data.ClientID = types.StringValue(res.GetClientId())
	data.PlexToken = types.StringValue(res.PlexToken)
	// Everything below has a Default and takes the API value verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.Promoted = types.BoolValue(res.GetPromoted())
	data.PolicyEngineMode = types.StringValue(string(res.GetPolicyEngineMode()))
	data.UserMatchingMode = types.StringValue(string(res.GetUserMatchingMode()))
	data.GroupMatchingMode = types.StringValue(string(res.GetGroupMatchingMode()))
	data.AllowFriends = types.BoolValue(res.GetAllowFriends())

	return diags
}

func (r *sourcePlexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourcePlexModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesPlexCreate(ctx).PlexSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourcePlexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data sourcePlexModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesPlexRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *sourcePlexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourcePlexModel
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

	res, hr, err := r.client.SourcesAPI.SourcesPlexUpdate(ctx, priorID.ValueString()).PlexSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourcePlexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourcePlexModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesPlexDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
