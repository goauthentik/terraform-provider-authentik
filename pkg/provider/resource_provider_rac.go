package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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
	_ resource.Resource                = &providerRACResource{}
	_ resource.ResourceWithConfigure   = &providerRACResource{}
	_ resource.ResourceWithImportState = &providerRACResource{}
)

func newProviderRACResource() resource.Resource {
	return &providerRACResource{}
}

type providerRACResource struct {
	resourceBase
}

type providerRACModel struct {
	ID                 types.String         `tfsdk:"id"`
	Name               types.String         `tfsdk:"name"`
	AuthenticationFlow types.String         `tfsdk:"authentication_flow"`
	AuthorizationFlow  types.String         `tfsdk:"authorization_flow"`
	PropertyMappings   types.List           `tfsdk:"property_mappings"`
	Settings           jsontypes.Normalized `tfsdk:"settings"`
	ConnectionExpiry   types.String         `tfsdk:"connection_expiry"`
}

func (r *providerRACResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_rac"
}

func (r *providerRACResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	settingsDefault := helpers.StringDefault("{}")
	connectionExpiryDefault := helpers.StringDefault("seconds=0")

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
			"authentication_flow": schema.StringAttribute{
				Optional: true,
			},
			"authorization_flow": schema.StringAttribute{
				Required: true,
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			// Note SDKv2 declared ValidateJSON here but *not* DiffSuppressJSON, unlike
			// authentik_property_mapping_provider_rac's settings. jsontypes.Normalized
			// brings both, so this gains semantic equality it did not have - which only ever
			// removes spurious diffs (`{"a": 1}` vs `{"a":1}`), never adds one. The protocol
			// type stays string, so the doc page is unchanged.
			"settings": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				Default:             settingsDefault,
				MarkdownDescription: helpers.Desc(helpers.JSONDescription, helpers.WithDefault(settingsDefault.Value())),
			},
			"connection_expiry": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             connectionExpiryDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(connectionExpiryDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
		},
	}
}

func (r *providerRACResource) toRequest(ctx context.Context, data *providerRACModel) (*api.RACProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	var settings map[string]any
	diags.Append(data.Settings.Unmarshal(&settings)...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.RACProviderRequest{
		Name:               data.Name.ValueString(),
		AuthorizationFlow:  data.AuthorizationFlow.ValueString(),
		PropertyMappings:   propertyMappings,
		ConnectionExpiry:   new(data.ConnectionExpiry.ValueString()),
		AuthenticationFlow: *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		Settings:           settings,
	}, diags
}

func (r *providerRACResource) fromAPI(ctx context.Context, data *providerRACModel, res *api.RACProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.AuthorizationFlow = types.StringValue(res.AuthorizationFlow)
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.ConnectionExpiry = types.StringValue(res.GetConnectionExpiry())

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	settingsBytes, err := json.Marshal(res.Settings)
	if err != nil {
		diags.AddError("Failed to encode RAC provider settings", err.Error())
		return diags
	}
	data.Settings = jsontypes.NewNormalizedValue(string(settingsBytes))

	return diags
}

func (r *providerRACResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerRACModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersRacCreate(ctx).RACProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerRACResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerRACModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersRacRetrieve(ctx, id).Execute()
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

func (r *providerRACResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerRACModel
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

	res, hr, err := r.client.ProvidersAPI.ProvidersRacUpdate(ctx, id).RACProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerRACResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerRACModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersRacDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
