package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &providerLDAPResource{}
	_ resource.ResourceWithConfigure   = &providerLDAPResource{}
	_ resource.ResourceWithImportState = &providerLDAPResource{}
)

func newProviderLDAPResource() resource.Resource {
	return &providerLDAPResource{}
}

type providerLDAPResource struct {
	resourceBase
}

type providerLDAPModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	BindFlow       types.String `tfsdk:"bind_flow"`
	UnbindFlow     types.String `tfsdk:"unbind_flow"`
	BaseDN         types.String `tfsdk:"base_dn"`
	Certificate    types.String `tfsdk:"certificate"`
	TLSServerName  types.String `tfsdk:"tls_server_name"`
	UIDStartNumber types.Int32  `tfsdk:"uid_start_number"`
	GIDStartNumber types.Int32  `tfsdk:"gid_start_number"`
	SearchMode     types.String `tfsdk:"search_mode"`
	BindMode       types.String `tfsdk:"bind_mode"`
	MFASupport     types.Bool   `tfsdk:"mfa_support"`
}

func (r *providerLDAPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_ldap"
}

func (r *providerLDAPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	uidStartNumberDefault := helpers.Int32Default(2000)
	gidStartNumberDefault := helpers.Int32Default(4000)
	searchModeDefault := helpers.StringDefault(string(api.LDAPAPIACCESSMODE_DIRECT))
	bindModeDefault := helpers.StringDefault(string(api.LDAPAPIACCESSMODE_DIRECT))
	mfaSupportDefault := helpers.BoolDefault(true)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Applications --- ",
		Attributes: map[string]schema.Attribute{
			// The providers are the int32-PK cluster: the PK is a stable server-generated
			// integer, stringified into state (see helpers.ParseInt32ID), so unlike the
			// slug-keyed sources it is safe to pin with UseStateForUnknown.
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"bind_flow": schema.StringAttribute{
				Required: true,
			},
			"unbind_flow": schema.StringAttribute{
				Required: true,
			},
			"base_dn": schema.StringAttribute{
				Required: true,
			},
			"certificate": schema.StringAttribute{
				Optional: true,
			},
			"tls_server_name": schema.StringAttribute{
				Optional: true,
			},
			"uid_start_number": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             uidStartNumberDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(uidStartNumberDefault.Value())),
			},
			"gid_start_number": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             gidStartNumberDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(gidStartNumberDefault.Value())),
			},
			// search_mode and bind_mode are api.LDAPAPIAccessMode on the wire but SDKv2
			// declared no validator and no enum description for either, and the doc page
			// shows them as bare Strings. Left that way deliberately - adding a OneOf here
			// would be a behaviour change and adding the prose would drift the docs.
			"search_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             searchModeDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(searchModeDefault.Value())),
			},
			"bind_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             bindModeDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(bindModeDefault.Value())),
			},
			"mfa_support": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             mfaSupportDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(mfaSupportDefault.Value())),
			},
		},
	}
}

func (r *providerLDAPResource) toRequest(data *providerLDAPModel) *api.LDAPProviderRequest {
	return &api.LDAPProviderRequest{
		Name:              data.Name.ValueString(),
		AuthorizationFlow: data.BindFlow.ValueString(),
		InvalidationFlow:  data.UnbindFlow.ValueString(),
		BaseDn:            new(data.BaseDN.ValueString()),
		UidStartNumber:    new(data.UIDStartNumber.ValueInt32()),
		GidStartNumber:    new(data.GIDStartNumber.ValueInt32()),
		SearchMode:        api.LDAPAPIAccessMode(data.SearchMode.ValueString()).Ptr(),
		BindMode:          api.LDAPAPIAccessMode(data.BindMode.ValueString()).Ptr(),
		MfaSupport:        new(data.MFASupport.ValueBool()),
		Certificate:       *api.NewNullableString(helpers.StringPtr(data.Certificate)),
		TlsServerName:     helpers.StringPtr(data.TLSServerName),
	}
}

func (r *providerLDAPResource) fromAPI(data *providerLDAPModel, res *api.LDAPProvider) {
	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.BindFlow = types.StringValue(res.AuthorizationFlow)
	data.UnbindFlow = types.StringValue(res.InvalidationFlow)
	// Nullable API field.
	data.Certificate = helpers.StringPtrOrNull(res.Certificate.Get())
	// No Default, so prior-aware.
	data.TLSServerName = helpers.StringOrNull(data.TLSServerName, res.GetTlsServerName())
	// Required in Terraform, so never null.
	data.BaseDN = types.StringValue(res.GetBaseDn())
	// The rest have Defaults and take the API value verbatim.
	data.UIDStartNumber = types.Int32Value(res.GetUidStartNumber())
	data.GIDStartNumber = types.Int32Value(res.GetGidStartNumber())
	data.SearchMode = types.StringValue(string(res.GetSearchMode()))
	data.BindMode = types.StringValue(string(res.GetBindMode()))
	data.MFASupport = types.BoolValue(res.GetMfaSupport())
}

func (r *providerLDAPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerLDAPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersLdapCreate(ctx).LDAPProviderRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerLDAPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersLdapRetrieve(ctx, id).Execute()
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

func (r *providerLDAPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerLDAPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The PK is stable and pinned with UseStateForUnknown, so the plan's id is already the
	// real one - unlike the slug-keyed sources, no read from req.State is needed here.
	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersLdapUpdate(ctx, id).LDAPProviderRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerLDAPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersLdapDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
