package provider

import (
	"context"
	"strconv"

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
	_ resource.Resource                = &providerRadiusResource{}
	_ resource.ResourceWithConfigure   = &providerRadiusResource{}
	_ resource.ResourceWithImportState = &providerRadiusResource{}
)

func newProviderRadiusResource() resource.Resource {
	return &providerRadiusResource{}
}

type providerRadiusResource struct {
	resourceBase
}

type providerRadiusModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	AuthorizationFlow types.String `tfsdk:"authorization_flow"`
	InvalidationFlow  types.String `tfsdk:"invalidation_flow"`
	PropertyMappings  types.List   `tfsdk:"property_mappings"`
	ClientNetworks    types.String `tfsdk:"client_networks"`
	SharedSecret      types.String `tfsdk:"shared_secret"`
	MFASupport        types.Bool   `tfsdk:"mfa_support"`
	Certificate       types.String `tfsdk:"certificate"`
}

func (r *providerRadiusResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_radius"
}

func (r *providerRadiusResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	clientNetworksDefault := helpers.StringDefault("0.0.0.0/0, ::/0")
	mfaSupportDefault := helpers.BoolDefault(true)

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
			"authorization_flow": schema.StringAttribute{
				Required: true,
			},
			"invalidation_flow": schema.StringAttribute{
				Required: true,
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"client_networks": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             clientNetworksDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(clientNetworksDefault.Value())),
			},
			// Unlike the other providers' secrets, shared_secret is returned by the API, so
			// this resource is not one of the 18 H4 cases and fromAPI writes it normally.
			"shared_secret": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"mfa_support": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             mfaSupportDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(mfaSupportDefault.Value())),
			},
			"certificate": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *providerRadiusResource) toRequest(ctx context.Context, data *providerRadiusModel) (*api.RadiusProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.RadiusProviderRequest{
		Name:              data.Name.ValueString(),
		AuthorizationFlow: data.AuthorizationFlow.ValueString(),
		InvalidationFlow:  data.InvalidationFlow.ValueString(),
		ClientNetworks:    new(data.ClientNetworks.ValueString()),
		SharedSecret:      new(data.SharedSecret.ValueString()),
		MfaSupport:        new(data.MFASupport.ValueBool()),
		PropertyMappings:  propertyMappings,
		Certificate:       *api.NewNullableString(helpers.StringPtr(data.Certificate)),
	}, diags
}

func (r *providerRadiusResource) fromAPI(ctx context.Context, data *providerRadiusModel, res *api.RadiusProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.AuthorizationFlow = types.StringValue(res.AuthorizationFlow)
	data.InvalidationFlow = types.StringValue(res.InvalidationFlow)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	// Nullable API field.
	data.Certificate = helpers.StringPtrOrNull(res.Certificate.Get())
	// Required, so never null.
	data.SharedSecret = types.StringValue(res.GetSharedSecret())
	// Have Defaults, so verbatim.
	data.ClientNetworks = types.StringValue(res.GetClientNetworks())
	data.MFASupport = types.BoolValue(res.GetMfaSupport())

	return diags
}

func (r *providerRadiusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerRadiusModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersRadiusCreate(ctx).RadiusProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerRadiusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerRadiusModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersRadiusRetrieve(ctx, id).Execute()
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

func (r *providerRadiusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerRadiusModel
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

	res, hr, err := r.client.ProvidersAPI.ProvidersRadiusUpdate(ctx, id).RadiusProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerRadiusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerRadiusModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersRadiusDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
