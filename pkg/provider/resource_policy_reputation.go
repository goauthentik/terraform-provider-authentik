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
	_ resource.Resource                = &policyReputationResource{}
	_ resource.ResourceWithConfigure   = &policyReputationResource{}
	_ resource.ResourceWithImportState = &policyReputationResource{}
)

func newPolicyReputationResource() resource.Resource {
	return &policyReputationResource{}
}

type policyReputationResource struct {
	resourceBase
}

type policyReputationModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ExecutionLogging types.Bool   `tfsdk:"execution_logging"`
	CheckIP          types.Bool   `tfsdk:"check_ip"`
	CheckUsername    types.Bool   `tfsdk:"check_username"`
	Threshold        types.Int32  `tfsdk:"threshold"`
}

func (r *policyReputationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_reputation"
}

func (r *policyReputationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)
	checkIPDefault := helpers.BoolDefault(true)
	checkUsernameDefault := helpers.BoolDefault(true)
	thresholdDefault := helpers.Int32Default(10)

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
			"check_ip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             checkIPDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(checkIPDefault.Value())),
			},
			"check_username": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             checkUsernameDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(checkUsernameDefault.Value())),
			},
			"threshold": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             thresholdDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(thresholdDefault.Value())),
			},
		},
	}
}

func (r *policyReputationResource) toRequest(data *policyReputationModel) *api.ReputationPolicyRequest {
	return &api.ReputationPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
		CheckIp:          new(data.CheckIP.ValueBool()),
		CheckUsername:    new(data.CheckUsername.ValueBool()),
		Threshold:        helpers.Int32Ptr(data.Threshold),
	}
}

func (r *policyReputationResource) fromAPI(data *policyReputationModel, res *api.ReputationPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	data.CheckIP = types.BoolValue(res.GetCheckIp())
	data.CheckUsername = types.BoolValue(res.GetCheckUsername())
	data.Threshold = types.Int32Value(res.GetThreshold())
}

func (r *policyReputationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyReputationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesReputationCreate(ctx).ReputationPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyReputationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyReputationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesReputationRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyReputationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyReputationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesReputationUpdate(ctx, data.ID.ValueString()).ReputationPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyReputationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyReputationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesReputationDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
