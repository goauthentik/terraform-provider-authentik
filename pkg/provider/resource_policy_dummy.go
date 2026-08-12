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
	_ resource.Resource                = &policyDummyResource{}
	_ resource.ResourceWithConfigure   = &policyDummyResource{}
	_ resource.ResourceWithImportState = &policyDummyResource{}
)

func newPolicyDummyResource() resource.Resource {
	return &policyDummyResource{}
}

type policyDummyResource struct {
	resourceBase
}

type policyDummyModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ExecutionLogging types.Bool   `tfsdk:"execution_logging"`
	Result           types.Bool   `tfsdk:"result"`
	WaitMin          types.Int32  `tfsdk:"wait_min"`
	WaitMax          types.Int32  `tfsdk:"wait_max"`
}

func (r *policyDummyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_dummy"
}

func (r *policyDummyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)
	resultDefault := helpers.BoolDefault(false)
	waitMinDefault := helpers.Int32Default(5)
	waitMaxDefault := helpers.Int32Default(30)

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
			"result": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             resultDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(resultDefault.Value())),
			},
			"wait_min": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             waitMinDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(waitMinDefault.Value())),
			},
			"wait_max": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             waitMaxDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(waitMaxDefault.Value())),
			},
		},
	}
}

func (r *policyDummyResource) toRequest(data *policyDummyModel) *api.DummyPolicyRequest {
	return &api.DummyPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
		Result:           new(data.Result.ValueBool()),
		WaitMin:          helpers.Int32Ptr(data.WaitMin),
		WaitMax:          helpers.Int32Ptr(data.WaitMax),
	}
}

// Every attribute here has a schema Default, so all of them take the API value verbatim
// (discovery #5) rather than going through the prior-aware helpers.
func (r *policyDummyResource) fromAPI(data *policyDummyModel, res *api.DummyPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	data.Result = types.BoolValue(res.GetResult())
	data.WaitMin = types.Int32Value(res.GetWaitMin())
	data.WaitMax = types.Int32Value(res.GetWaitMax())
}

func (r *policyDummyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyDummyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesDummyCreate(ctx).DummyPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyDummyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyDummyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesDummyRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyDummyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyDummyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesDummyUpdate(ctx, data.ID.ValueString()).DummyPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyDummyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyDummyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesDummyDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
