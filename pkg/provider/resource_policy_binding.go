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
	_ resource.Resource                = &policyBindingResource{}
	_ resource.ResourceWithConfigure   = &policyBindingResource{}
	_ resource.ResourceWithImportState = &policyBindingResource{}
)

func newPolicyBindingResource() resource.Resource {
	return &policyBindingResource{}
}

type policyBindingResource struct {
	resourceBase
}

// policyBindingModel is the only policy without a name or execution_logging: a binding is
// an edge between a target and one of policy/user/group rather than a policy itself.
type policyBindingModel struct {
	ID            types.String `tfsdk:"id"`
	Target        types.String `tfsdk:"target"`
	Policy        types.String `tfsdk:"policy"`
	User          types.Int32  `tfsdk:"user"`
	Group         types.String `tfsdk:"group"`
	Order         types.Int32  `tfsdk:"order"`
	Negate        types.Bool   `tfsdk:"negate"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Timeout       types.Int32  `tfsdk:"timeout"`
	FailureResult types.Bool   `tfsdk:"failure_result"`
}

func (r *policyBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_binding"
}

func (r *policyBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	negateDefault := helpers.BoolDefault(false)
	enabledDefault := helpers.BoolDefault(true)
	timeoutDefault := helpers.Int32Default(30)
	failureResultDefault := helpers.BoolDefault(false)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the object this binding should apply to",
			},
			"policy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "UUID of the policy",
			},
			"user": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "PK of the user",
			},
			"group": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "UUID of the group",
			},
			"order": schema.Int32Attribute{
				Required: true,
			},
			"negate": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             negateDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(negateDefault.Value())),
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enabledDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enabledDefault.Value())),
			},
			"timeout": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             timeoutDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(timeoutDefault.Value())),
			},
			"failure_result": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             failureResultDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(failureResultDefault.Value())),
			},
		},
	}
}

func (r *policyBindingResource) toRequest(data *policyBindingModel) *api.PolicyBindingRequest {
	return &api.PolicyBindingRequest{
		Target:        data.Target.ValueString(),
		Order:         data.Order.ValueInt32(),
		Negate:        new(data.Negate.ValueBool()),
		Enabled:       new(data.Enabled.ValueBool()),
		Timeout:       new(data.Timeout.ValueInt32()),
		FailureResult: new(data.FailureResult.ValueBool()),
		// policy/user/group are all genuinely nullable API fields and mutually exclusive in
		// practice, so a null config must serialise as null rather than as "" or 0.
		Policy: *api.NewNullableString(helpers.StringPtr(data.Policy)),
		User:   *api.NewNullableInt32(helpers.Int32Ptr(data.User)),
		Group:  *api.NewNullableString(helpers.StringPtr(data.Group)),
	}
}

func (r *policyBindingResource) fromAPI(data *policyBindingModel, res *api.PolicyBinding) {
	data.ID = types.StringValue(res.Pk)
	data.Target = types.StringValue(res.Target)
	// order is Required and a plain int32 on the wire, so never null.
	data.Order = types.Int32Value(res.Order)
	// Nullable API fields, so the Ptr variants are correct here rather than the prior-aware
	// helpers - see discovery #3.
	data.Policy = helpers.StringPtrOrNull(res.Policy.Get())
	data.User = helpers.Int32PtrOrNull(res.User.Get())
	data.Group = helpers.StringPtrOrNull(res.Group.Get())
	// The remaining four have Defaults and take the API value verbatim.
	data.Negate = types.BoolValue(res.GetNegate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.Timeout = types.Int32Value(res.GetTimeout())
	data.FailureResult = types.BoolValue(res.GetFailureResult())
}

func (r *policyBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyBindingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesBindingsCreate(ctx).PolicyBindingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyBindingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesBindingsRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyBindingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesBindingsUpdate(ctx, data.ID.ValueString()).PolicyBindingRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyBindingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesBindingsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
