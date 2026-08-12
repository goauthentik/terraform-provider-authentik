package provider

import (
	"context"

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
	_ resource.Resource                = &sourceKerberosResource{}
	_ resource.ResourceWithConfigure   = &sourceKerberosResource{}
	_ resource.ResourceWithImportState = &sourceKerberosResource{}
)

func newSourceKerberosResource() resource.Resource {
	return &sourceKerberosResource{}
}

type sourceKerberosResource struct {
	resourceBase
}

type sourceKerberosModel struct {
	ID                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	UUID                                types.String `tfsdk:"uuid"`
	Slug                                types.String `tfsdk:"slug"`
	UserPathTemplate                    types.String `tfsdk:"user_path_template"`
	AuthenticationFlow                  types.String `tfsdk:"authentication_flow"`
	EnrollmentFlow                      types.String `tfsdk:"enrollment_flow"`
	Enabled                             types.Bool   `tfsdk:"enabled"`
	PolicyEngineMode                    types.String `tfsdk:"policy_engine_mode"`
	UserMatchingMode                    types.String `tfsdk:"user_matching_mode"`
	GroupMatchingMode                   types.String `tfsdk:"group_matching_mode"`
	Realm                               types.String `tfsdk:"realm"`
	Krb5Conf                            types.String `tfsdk:"krb5_conf"`
	SyncUsers                           types.Bool   `tfsdk:"sync_users"`
	SyncUsersPassword                   types.Bool   `tfsdk:"sync_users_password"`
	SyncPrincipal                       types.String `tfsdk:"sync_principal"`
	SyncPassword                        types.String `tfsdk:"sync_password"`
	SyncKeytab                          types.String `tfsdk:"sync_keytab"`
	SyncCcache                          types.String `tfsdk:"sync_ccache"`
	SpnegoServerName                    types.String `tfsdk:"spnego_server_name"`
	SpnegoKeytab                        types.String `tfsdk:"spnego_keytab"`
	SpnegoCcache                        types.String `tfsdk:"spnego_ccache"`
	PasswordLoginUpdateInternalPassword types.Bool   `tfsdk:"password_login_update_internal_password"`
	SyncOutgoingTriggerMode             types.String `tfsdk:"sync_outgoing_trigger_mode"`
}

func (r *sourceKerberosResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_kerberos"
}

func (r *sourceKerberosResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	userMatchingModeDefault := helpers.StringDefault(string(api.USERMATCHINGMODEENUM_IDENTIFIER))
	groupMatchingModeDefault := helpers.StringDefault(string(api.GROUPMATCHINGMODEENUM_IDENTIFIER))
	syncUsersDefault := helpers.BoolDefault(true)
	syncUsersPasswordDefault := helpers.BoolDefault(true)
	passwordLoginUpdateInternalPasswordDefault := helpers.BoolDefault(false)
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
			"realm": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Kerberos realm",
			},
			"krb5_conf": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Custom krb5.conf to use. Uses the system one by default",
			},
			"sync_users": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncUsersDefault,
				MarkdownDescription: helpers.Desc("Sync users from Kerberos into authentik", helpers.WithDefault(syncUsersDefault.Value())),
			},
			"sync_users_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             syncUsersPasswordDefault,
				MarkdownDescription: helpers.Desc("When a user changes their password, sync it back to Kerberos", helpers.WithDefault(syncUsersPasswordDefault.Value())),
			},
			"sync_principal": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Principal to authenticate to kadmin for sync.",
			},
			// sync_password, sync_keytab and spnego_keytab are three of the 18 H4
			// attributes: the API never returns them, so fromAPI must not assign them.
			"sync_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password to authenticate to kadmin for sync",
			},
			"sync_keytab": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Keytab to authenticate to kadmin for sync. Must be base64-encoded or in the form TYPE:residual",
			},
			"sync_ccache": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Credentials cache to authenticate to kadmin for sync. Must be in the form TYPE:residual",
			},
			"spnego_server_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Force the use of a specific server name for SPNEGO",
			},
			"spnego_keytab": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "SPNEGO keytab base64-encoded or path to keytab in the form FILE:path",
			},
			"spnego_ccache": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Credential cache to use for SPNEGO in form type:residual",
			},
			"password_login_update_internal_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             passwordLoginUpdateInternalPasswordDefault,
				MarkdownDescription: helpers.Desc("If enabled, the authentik-stored password will be updated upon login with the Kerberos password backend", helpers.WithDefault(passwordLoginUpdateInternalPasswordDefault.Value())),
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

func (r *sourceKerberosResource) toRequest(data *sourceKerberosModel) *api.KerberosSourceRequest {
	return &api.KerberosSourceRequest{
		Name:               data.Name.ValueString(),
		Slug:               data.Slug.ValueString(),
		Enabled:            new(data.Enabled.ValueBool()),
		UserPathTemplate:   new(data.UserPathTemplate.ValueString()),
		PolicyEngineMode:   api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		UserMatchingMode:   api.UserMatchingModeEnum(data.UserMatchingMode.ValueString()).Ptr(),
		GroupMatchingMode:  api.GroupMatchingModeEnum(data.GroupMatchingMode.ValueString()).Ptr(),
		AuthenticationFlow: *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		EnrollmentFlow:     *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		Realm:              data.Realm.ValueString(),
		// SDKv2 built all of these with new(d.Get(...).(string)), i.e. always a pointer and
		// "" when unset, so they keep StringPtrEmpty semantics rather than being omitted.
		Krb5Conf:                            new(data.Krb5Conf.ValueString()),
		SyncPrincipal:                       new(data.SyncPrincipal.ValueString()),
		SyncPassword:                        new(data.SyncPassword.ValueString()),
		SyncKeytab:                          new(data.SyncKeytab.ValueString()),
		SyncCcache:                          new(data.SyncCcache.ValueString()),
		SpnegoServerName:                    new(data.SpnegoServerName.ValueString()),
		SpnegoKeytab:                        new(data.SpnegoKeytab.ValueString()),
		SpnegoCcache:                        new(data.SpnegoCcache.ValueString()),
		SyncUsers:                           new(data.SyncUsers.ValueBool()),
		SyncUsersPassword:                   new(data.SyncUsersPassword.ValueBool()),
		PasswordLoginUpdateInternalPassword: new(data.PasswordLoginUpdateInternalPassword.ValueBool()),
		SyncOutgoingTriggerMode:             api.SyncOutgoingTriggerModeEnum(data.SyncOutgoingTriggerMode.ValueString()).Ptr(),
	}
}

// fromAPI deliberately never assigns SyncPassword, SyncKeytab or SpnegoKeytab: this is one
// of the 18 H4 resources and the API returns none of the three. Callers seed data from
// req.Plan or req.State first, so the omissions preserve the configured secrets - see
// stageCaptchaResource.fromAPI for the full explanation.
func (r *sourceKerberosResource) fromAPI(data *sourceKerberosModel, res *api.KerberosSource) {
	// H3: id is the slug, not the Pk.
	data.ID = types.StringValue(res.Slug)
	data.Name = types.StringValue(res.Name)
	data.Slug = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)
	data.Realm = types.StringValue(res.Realm)

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.EnrollmentFlow = helpers.StringPtrOrNull(res.EnrollmentFlow.Get())

	// No Default on these five, so they keep the prior-aware helper.
	data.Krb5Conf = helpers.StringOrNull(data.Krb5Conf, res.GetKrb5Conf())
	data.SyncPrincipal = helpers.StringOrNull(data.SyncPrincipal, res.GetSyncPrincipal())
	data.SyncCcache = helpers.StringOrNull(data.SyncCcache, res.GetSyncCcache())
	data.SpnegoServerName = helpers.StringOrNull(data.SpnegoServerName, res.GetSpnegoServerName())
	data.SpnegoCcache = helpers.StringOrNull(data.SpnegoCcache, res.GetSpnegoCcache())

	// Everything below has a Default and takes the API value verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.PolicyEngineMode = types.StringValue(string(res.GetPolicyEngineMode()))
	data.UserMatchingMode = types.StringValue(string(res.GetUserMatchingMode()))
	// SDKv2's Read set group_matching_mode from res.UserMatchingMode - a copy-paste slip
	// that no other source has. Under SDKv2 it silently produced a permanent diff whenever
	// the two modes differed; in the framework it would be a hard "provider produced
	// inconsistent result" error, since the plan carries the configured value and state
	// would get the other one. Fixed here to read GroupMatchingMode.
	data.GroupMatchingMode = types.StringValue(string(res.GetGroupMatchingMode()))
	data.SyncUsers = types.BoolValue(res.GetSyncUsers())
	data.SyncUsersPassword = types.BoolValue(res.GetSyncUsersPassword())
	data.PasswordLoginUpdateInternalPassword = types.BoolValue(res.GetPasswordLoginUpdateInternalPassword())
	data.SyncOutgoingTriggerMode = types.StringValue(string(res.GetSyncOutgoingTriggerMode()))
}

func (r *sourceKerberosResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceKerberosModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesKerberosCreate(ctx).KerberosSourceRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceKerberosResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state carries the three keytab/password secrets - see fromAPI.
	var data sourceKerberosModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesKerberosRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceKerberosResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceKerberosModel
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

	res, hr, err := r.client.SourcesAPI.SourcesKerberosUpdate(ctx, priorID.ValueString()).KerberosSourceRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceKerberosResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceKerberosModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesKerberosDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
