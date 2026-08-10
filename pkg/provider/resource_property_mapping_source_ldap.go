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
	_ resource.Resource                = &propertyMappingSourceLDAPResource{}
	_ resource.ResourceWithConfigure   = &propertyMappingSourceLDAPResource{}
	_ resource.ResourceWithImportState = &propertyMappingSourceLDAPResource{}
)

func newPropertyMappingSourceLDAPResource() resource.Resource {
	return &propertyMappingSourceLDAPResource{}
}

type propertyMappingSourceLDAPResource struct {
	resourceBase
}

type propertyMappingSourceLDAPModel struct {
	ID         types.String            `tfsdk:"id"`
	Name       types.String            `tfsdk:"name"`
	Expression helpers.ExpressionValue `tfsdk:"expression"`
}

func (r *propertyMappingSourceLDAPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_property_mapping_source_ldap"
}

func (r *propertyMappingSourceLDAPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- Manage LDAP Source Property mappings",
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
			// helpers.ExpressionType replaces SDKv2's DiffSuppressExpression: authentik
			// returns expressions with trailing newlines stripped, so a heredoc config
			// value ending in "\n" would otherwise show a permanent diff. It stays a
			// protocol `string`, so neither the doc page nor the state layout changes.
			"expression": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Required:   true,
			},
		},
	}
}

func (r *propertyMappingSourceLDAPResource) toRequest(data *propertyMappingSourceLDAPModel) *api.LDAPSourcePropertyMappingRequest {
	return &api.LDAPSourcePropertyMappingRequest{
		Name:       data.Name.ValueString(),
		Expression: data.Expression.ValueString(),
	}
}

// fromAPI needs no prior-aware helper for expression: it is Required, so it can never be
// null, and semantic equality handles the trailing-newline difference on its own.
func (r *propertyMappingSourceLDAPResource) fromAPI(data *propertyMappingSourceLDAPModel, res *api.LDAPSourcePropertyMapping) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Expression = helpers.NewExpressionValue(res.Expression)
}

func (r *propertyMappingSourceLDAPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data propertyMappingSourceLDAPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsSourceLdapCreate(ctx).LDAPSourcePropertyMappingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingSourceLDAPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data propertyMappingSourceLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsSourceLdapRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *propertyMappingSourceLDAPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data propertyMappingSourceLDAPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PropertymappingsAPI.PropertymappingsSourceLdapUpdate(ctx, data.ID.ValueString()).LDAPSourcePropertyMappingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *propertyMappingSourceLDAPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data propertyMappingSourceLDAPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PropertymappingsAPI.PropertymappingsSourceLdapDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
