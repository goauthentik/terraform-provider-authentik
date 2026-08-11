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
	_ resource.Resource                = &sourceLDAPResource{}
	_ resource.ResourceWithConfigure   = &sourceLDAPResource{}
	_ resource.ResourceWithImportState = &sourceLDAPResource{}
)

func newSourceLDAPResource() resource.Resource {
	return &sourceLDAPResource{}
}

type sourceLDAPResource struct {
	resourceBase
}

// sourceLDAPModel has no authentication_flow/enrollment_flow/policy_engine_mode/
// user_matching_mode, unlike every other source: LDAP is sync-only rather than an
// interactive login source.
type sourceLDAPModel struct {
	ID                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	UUID                                types.String `tfsdk:"uuid"`
	Slug                                types.String `tfsdk:"slug"`
	UserPathTemplate                    types.String `tfsdk:"user_path_template"`
	Enabled                             types.Bool   `tfsdk:"enabled"`
	ServerURI                           types.String `tfsdk:"server_uri"`
	BindCN                              types.String `tfsdk:"bind_cn"`
	BindPassword                        types.String `tfsdk:"bind_password"`
	StartTLS                            types.Bool   `tfsdk:"start_tls"`
	SNI                                 types.Bool   `tfsdk:"sni"`
	BaseDN                              types.String `tfsdk:"base_dn"`
	AdditionalUserDN                    types.String `tfsdk:"additional_user_dn"`
	AdditionalGroupDN                   types.String `tfsdk:"additional_group_dn"`
	UserObjectFilter                    types.String `tfsdk:"user_object_filter"`
	UserMembershipAttribute             types.String `tfsdk:"user_membership_attribute"`
	GroupObjectFilter                   types.String `tfsdk:"group_object_filter"`
	GroupMembershipField                types.String `tfsdk:"group_membership_field"`
	ObjectUniquenessField               types.String `tfsdk:"object_uniqueness_field"`
	LookupGroupsFromUser                types.Bool   `tfsdk:"lookup_groups_from_user"`
	SyncUsers                           types.Bool   `tfsdk:"sync_users"`
	SyncUsersPassword                   types.Bool   `tfsdk:"sync_users_password"`
	SyncGroups                          types.Bool   `tfsdk:"sync_groups"`
	SyncParentGroup                     types.String `tfsdk:"sync_parent_group"`
	PasswordLoginUpdateInternalPassword types.Bool   `tfsdk:"password_login_update_internal_password"`
	DeleteNotFoundObjects               types.Bool   `tfsdk:"delete_not_found_objects"`
	PropertyMappings                    types.List   `tfsdk:"property_mappings"`
	PropertyMappingsGroup               types.List   `tfsdk:"property_mappings_group"`
	SyncOutgoingTriggerMode             types.String `tfsdk:"sync_outgoing_trigger_mode"`
}

func (r *sourceLDAPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_ldap"
}

func (r *sourceLDAPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	startTLSDefault := helpers.BoolDefault(true)
	sniDefault := helpers.BoolDefault(false)
	additionalUserDNDefault := helpers.StringDefault("")
	additionalGroupDNDefault := helpers.StringDefault("")
	userObjectFilterDefault := helpers.StringDefault("(objectClass=person)")
	userMembershipAttributeDefault := helpers.StringDefault("distinguishedName")
	groupObjectFilterDefault := helpers.StringDefault("(objectClass=group)")
	groupMembershipFieldDefault := helpers.StringDefault("member")
	objectUniquenessFieldDefault := helpers.StringDefault("objectSid")
	lookupGroupsFromUserDefault := helpers.BoolDefault(true)
	syncUsersDefault := helpers.BoolDefault(true)
	syncUsersPasswordDefault := helpers.BoolDefault(true)
	syncGroupsDefault := helpers.BoolDefault(true)
	passwordLoginUpdateInternalPasswordDefault := helpers.BoolDefault(false)
	deleteNotFoundObjectsDefault := helpers.BoolDefault(false)
	syncOutgoingTriggerModeDefault := helpers.StringDefault(string(api.SYNCOUTGOINGTRIGGERMODEENUM_DEFERRED_END))

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
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enabledDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enabledDefault.Value())),
			},
			"server_uri": schema.StringAttribute{
				Required: true,
			},
			"bind_cn": schema.StringAttribute{
				Required: true,
			},
			// One of the 18 H4 attributes: the API never returns bind_password.
			"bind_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"start_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             startTLSDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(startTLSDefault.Value())),
			},
			"sni": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             sniDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(sniDefault.Value())),
			},
			"base_dn": schema.StringAttribute{
				Required: true,
			},
			"additional_user_dn": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             additionalUserDNDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(additionalUserDNDefault.Value())),
			},
			"additional_group_dn": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             additionalGroupDNDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(additionalGroupDNDefault.Value())),
			},
			"user_object_filter": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userObjectFilterDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(userObjectFilterDefault.Value())),
			},
			"user_membership_attribute": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userMembershipAttributeDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(userMembershipAttributeDefault.Value())),
			},
			"group_object_filter": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             groupObjectFilterDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(groupObjectFilterDefault.Value())),
			},
			"group_membership_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             groupMembershipFieldDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(groupMembershipFieldDefault.Value())),
			},
			"object_uniqueness_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             objectUniquenessFieldDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(objectUniquenessFieldDefault.Value())),
			},
			"lookup_groups_from_user": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             lookupGroupsFromUserDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(lookupGroupsFromUserDefault.Value())),
			},
			"sync_users": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncUsersDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(syncUsersDefault.Value())),
			},
			"sync_users_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncUsersPasswordDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(syncUsersPasswordDefault.Value())),
			},
			"sync_groups": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncGroupsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(syncGroupsDefault.Value())),
			},
			"sync_parent_group": schema.StringAttribute{
				Optional: true,
			},
			"password_login_update_internal_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             passwordLoginUpdateInternalPasswordDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(passwordLoginUpdateInternalPasswordDefault.Value())),
			},
			"delete_not_found_objects": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             deleteNotFoundObjectsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(deleteNotFoundObjectsDefault.Value())),
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"property_mappings_group": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"sync_outgoing_trigger_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncOutgoingTriggerModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedSyncOutgoingTriggerModeEnumEnumValues), helpers.WithDefault(syncOutgoingTriggerModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedSyncOutgoingTriggerModeEnumEnumValues),
				},
			},
		},
	}
}

func (r *sourceLDAPResource) toRequest(ctx context.Context, data *sourceLDAPModel) (*api.LDAPSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.LDAPSourceRequest{
		Name:                                data.Name.ValueString(),
		Slug:                                data.Slug.ValueString(),
		Enabled:                             new(data.Enabled.ValueBool()),
		UserPathTemplate:                    new(data.UserPathTemplate.ValueString()),
		BaseDn:                              data.BaseDN.ValueString(),
		ServerUri:                           data.ServerURI.ValueString(),
		BindCn:                              new(data.BindCN.ValueString()),
		BindPassword:                        new(data.BindPassword.ValueString()),
		StartTls:                            new(data.StartTLS.ValueBool()),
		Sni:                                 new(data.SNI.ValueBool()),
		AdditionalUserDn:                    new(data.AdditionalUserDN.ValueString()),
		AdditionalGroupDn:                   new(data.AdditionalGroupDN.ValueString()),
		UserObjectFilter:                    new(data.UserObjectFilter.ValueString()),
		UserMembershipAttribute:             new(data.UserMembershipAttribute.ValueString()),
		GroupObjectFilter:                   new(data.GroupObjectFilter.ValueString()),
		GroupMembershipField:                new(data.GroupMembershipField.ValueString()),
		ObjectUniquenessField:               new(data.ObjectUniquenessField.ValueString()),
		SyncUsers:                           new(data.SyncUsers.ValueBool()),
		SyncUsersPassword:                   new(data.SyncUsersPassword.ValueBool()),
		SyncGroups:                          new(data.SyncGroups.ValueBool()),
		SyncParentGroup:                     *api.NewNullableString(helpers.StringPtr(data.SyncParentGroup)),
		PasswordLoginUpdateInternalPassword: new(data.PasswordLoginUpdateInternalPassword.ValueBool()),
		DeleteNotFoundObjects:               new(data.DeleteNotFoundObjects.ValueBool()),
		LookupGroupsFromUser:                new(data.LookupGroupsFromUser.ValueBool()),
		UserPropertyMappings:                propertyMappings,
		GroupPropertyMappings:               propertyMappingsGroup,
		SyncOutgoingTriggerMode:             api.SyncOutgoingTriggerModeEnum(data.SyncOutgoingTriggerMode.ValueString()).Ptr(),
	}, diags
}

// fromAPI deliberately never assigns BindPassword: this is one of the 18 H4 resources and the
// API does not return it. Callers seed data from req.Plan or req.State first, so the omission
// preserves the configured secret - see stageCaptchaResource.fromAPI. bind_password is
// Required, so nulling it would fail the apply rather than merely lose the value.
func (r *sourceLDAPResource) fromAPI(ctx context.Context, data *sourceLDAPModel, res *api.LDAPSource) diag.Diagnostics {
	var diags diag.Diagnostics

	// H3: id is the slug, not the Pk.
	data.ID = types.StringValue(res.Slug)
	data.Name = types.StringValue(res.Name)
	data.Slug = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)
	data.BaseDN = types.StringValue(res.BaseDn)
	data.ServerURI = types.StringValue(res.ServerUri)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.UserPropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	propertyMappingsGroup, d := helpers.MergeStringList(ctx, data.PropertyMappingsGroup, res.GroupPropertyMappings)
	diags.Append(d...)
	data.PropertyMappingsGroup = propertyMappingsGroup

	// Nullable API field.
	data.SyncParentGroup = helpers.StringPtrOrNull(res.SyncParentGroup.Get())
	// Required, so never null.
	data.BindCN = types.StringValue(res.GetBindCn())
	// Everything below has a Default and takes the API value verbatim. additional_user_dn
	// and additional_group_dn are the ones that make this matter: both default to "", which
	// is exactly what the API returns when unset, so the prior-aware helper would import
	// them as null and then plan a spurious update.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.StartTLS = types.BoolValue(res.GetStartTls())
	data.SNI = types.BoolValue(res.GetSni())
	data.AdditionalUserDN = types.StringValue(res.GetAdditionalUserDn())
	data.AdditionalGroupDN = types.StringValue(res.GetAdditionalGroupDn())
	data.UserObjectFilter = types.StringValue(res.GetUserObjectFilter())
	data.UserMembershipAttribute = types.StringValue(res.GetUserMembershipAttribute())
	data.GroupObjectFilter = types.StringValue(res.GetGroupObjectFilter())
	data.GroupMembershipField = types.StringValue(res.GetGroupMembershipField())
	data.ObjectUniquenessField = types.StringValue(res.GetObjectUniquenessField())
	data.LookupGroupsFromUser = types.BoolValue(res.GetLookupGroupsFromUser())
	data.SyncUsers = types.BoolValue(res.GetSyncUsers())
	data.SyncUsersPassword = types.BoolValue(res.GetSyncUsersPassword())
	data.SyncGroups = types.BoolValue(res.GetSyncGroups())
	data.PasswordLoginUpdateInternalPassword = types.BoolValue(res.GetPasswordLoginUpdateInternalPassword())
	data.DeleteNotFoundObjects = types.BoolValue(res.GetDeleteNotFoundObjects())
	data.SyncOutgoingTriggerMode = types.StringValue(string(res.GetSyncOutgoingTriggerMode()))

	return diags
}

func (r *sourceLDAPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceLDAPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesLdapCreate(ctx).LDAPSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceLDAPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries bind_password across a Read - see fromAPI.
	var data sourceLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesLdapRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *sourceLDAPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceLDAPModel
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

	res, hr, err := r.client.SourcesAPI.SourcesLdapUpdate(ctx, priorID.ValueString()).LDAPSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceLDAPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesLdapDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
