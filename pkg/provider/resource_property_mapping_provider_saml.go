package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &propertyMappingProviderSAMLResource{}
	_ resource.ResourceWithConfigure   = &propertyMappingProviderSAMLResource{}
	_ resource.ResourceWithImportState = &propertyMappingProviderSAMLResource{}
)

func newPropertyMappingProviderSAMLResource() resource.Resource {
	return &propertyMappingProviderSAMLResource{}
}

type propertyMappingProviderSAMLResource struct {
	resourceBase
}

type propertyMappingProviderSAMLModel struct {
	ID           types.String            `tfsdk:"id"`
	Name         types.String            `tfsdk:"name"`
	SamlName     types.String            `tfsdk:"saml_name"`
	FriendlyName types.String            `tfsdk:"friendly_name"`
	Expression   helpers.ExpressionValue `tfsdk:"expression"`
}

func (r *propertyMappingProviderSAMLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_mapping_provider_saml"
}

func (r *propertyMappingProviderSAMLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- Manage SAML Provider Property mappings",
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
			"saml_name": schema.StringAttribute{
				Required: true,
			},
			"friendly_name": schema.StringAttribute{
				Optional: true,
			},
			// See the note in resource_property_mapping_source_ldap.go on ExpressionType.
			"expression": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Required:   true,
			},
		},
	}
}

func (r *propertyMappingProviderSAMLResource) toRequest(data *propertyMappingProviderSAMLModel) *api.SAMLPropertyMappingRequest {
	return &api.SAMLPropertyMappingRequest{
		Name:       data.Name.ValueString(),
		SamlName:   data.SamlName.ValueString(),
		Expression: data.Expression.ValueString(),
		// friendly_name is a genuinely nullable API field (NullableString), so null config
		// stays null on the wire rather than becoming "".
		FriendlyName: *api.NewNullableString(helpers.StringPtr(data.FriendlyName)),
	}
}

func (r *propertyMappingProviderSAMLResource) fromAPI(data *propertyMappingProviderSAMLModel, res *api.SAMLPropertyMapping) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.SamlName = types.StringValue(res.SamlName)
	data.Expression = helpers.NewExpressionValue(res.Expression)
	// Nullable API field, so StringPtrOrNull is correct here rather than the prior-aware
	// StringOrNull - see discovery #3.
	data.FriendlyName = helpers.StringPtrOrNull(res.FriendlyName.Get())
}

func (r *propertyMappingProviderSAMLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data propertyMappingProviderSAMLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderSamlCreate(ctx).SAMLPropertyMappingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingProviderSAMLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data propertyMappingProviderSAMLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderSamlRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *propertyMappingProviderSAMLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data propertyMappingProviderSAMLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderSamlUpdate(ctx, data.ID.ValueString()).SAMLPropertyMappingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingProviderSAMLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data propertyMappingProviderSAMLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PropertymappingsAPI.PropertymappingsProviderSamlDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
