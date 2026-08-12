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
	_ resource.Resource                = &policyExpressionResource{}
	_ resource.ResourceWithConfigure   = &policyExpressionResource{}
	_ resource.ResourceWithImportState = &policyExpressionResource{}
)

func newPolicyExpressionResource() resource.Resource {
	return &policyExpressionResource{}
}

type policyExpressionResource struct {
	resourceBase
}

type policyExpressionModel struct {
	ID               types.String            `tfsdk:"id"`
	Name             types.String            `tfsdk:"name"`
	ExecutionLogging types.Bool              `tfsdk:"execution_logging"`
	Expression       helpers.ExpressionValue `tfsdk:"expression"`
}

func (r *policyExpressionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_expression"
}

func (r *policyExpressionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- ",
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
			"execution_logging": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             executionLoggingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(executionLoggingDefault.Value())),
			},
			// The 17th of the 18 DiffSuppressExpression attributes; see
			// resource_property_mapping_source_ldap.go and the registry sweep in
			// expression_registry_test.go, which covers this one too.
			"expression": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Required:   true,
			},
		},
	}
}

func (r *policyExpressionResource) toRequest(data *policyExpressionModel) *api.ExpressionPolicyRequest {
	return &api.ExpressionPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
		Expression:       data.Expression.ValueString(),
	}
}

func (r *policyExpressionResource) fromAPI(data *policyExpressionModel, res *api.ExpressionPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	// expression is Required, so it can never be null and needs no prior-aware helper;
	// semantic equality handles the trailing-newline difference on its own.
	data.Expression = helpers.NewExpressionValue(res.Expression)
}

func (r *policyExpressionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyExpressionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesExpressionCreate(ctx).ExpressionPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyExpressionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyExpressionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesExpressionRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyExpressionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyExpressionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesExpressionUpdate(ctx, data.ID.ValueString()).ExpressionPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyExpressionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyExpressionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesExpressionDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
