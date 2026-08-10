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
	_ resource.Resource                = &propertyMappingProviderRACResource{}
	_ resource.ResourceWithConfigure   = &propertyMappingProviderRACResource{}
	_ resource.ResourceWithImportState = &propertyMappingProviderRACResource{}
)

func newPropertyMappingProviderRACResource() resource.Resource {
	return &propertyMappingProviderRACResource{}
}

type propertyMappingProviderRACResource struct {
	resourceBase
}

type propertyMappingProviderRACModel struct {
	ID         types.String            `tfsdk:"id"`
	Name       types.String            `tfsdk:"name"`
	Expression helpers.ExpressionValue `tfsdk:"expression"`
	Settings   jsontypes.Normalized    `tfsdk:"settings"`
}

func (r *propertyMappingProviderRACResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_mapping_provider_rac"
}

func (r *propertyMappingProviderRACResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	settingsDefault := helpers.StringDefault("{}")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- Manage RAC Provider Property mappings",
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
			// The only property mapping whose expression is Optional rather than Required,
			// so it is also the only one needing ExpressionOrNull on the read side. See the
			// note in resource_property_mapping_source_ldap.go on ExpressionType itself.
			"expression": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Optional:   true,
			},
			// jsontypes.Normalized replaces SDKv2's DiffSuppressJSON + ValidateJSON pair:
			// it validates in ValidateAttribute and compares semantically, so {"a": 1} and
			// {"a":1} produce an empty plan. The JSONDescription prose is kept even though
			// ValidateJSON is gone, or the doc page changes.
			"settings": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				Optional:            true,
				Computed:            true,
				Default:             settingsDefault,
				MarkdownDescription: helpers.Desc(helpers.JSONDescription, helpers.WithDefault(settingsDefault.Value())),
			},
		},
	}
}

func (r *propertyMappingProviderRACResource) toRequest(data *propertyMappingProviderRACModel) (*api.RACPropertyMappingRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	var settings map[string]any
	diags.Append(data.Settings.Unmarshal(&settings)...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.RACPropertyMappingRequest{
		Name:           data.Name.ValueString(),
		Expression:     helpers.ExpressionPtr(data.Expression),
		StaticSettings: settings,
	}, diags
}

func (r *propertyMappingProviderRACResource) fromAPI(data *propertyMappingProviderRACModel, res *api.RACPropertyMapping) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	// Optional with no Default, so prior-aware: semantic equality cannot bridge null to ""
	// because value_semantic_equality.go returns early on a null prior value.
	data.Expression = helpers.ExpressionOrNull(data.Expression, res.GetExpression())

	settingsBytes, err := json.Marshal(res.StaticSettings)
	if err != nil {
		diags.AddError("Failed to encode RAC property mapping settings", err.Error())
		return diags
	}
	data.Settings = jsontypes.NewNormalizedValue(string(settingsBytes))

	return diags
}

func (r *propertyMappingProviderRACResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data propertyMappingProviderRACModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderRacCreate(ctx).RACPropertyMappingRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingProviderRACResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data propertyMappingProviderRACModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderRacRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingProviderRACResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data propertyMappingProviderRACModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(&data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderRacUpdate(ctx, data.ID.ValueString()).RACPropertyMappingRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(&data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingProviderRACResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data propertyMappingProviderRACModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderRacDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
