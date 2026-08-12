package provider

import (
	"context"

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
	_ resource.Resource                = &policyEventMatcherResource{}
	_ resource.ResourceWithConfigure   = &policyEventMatcherResource{}
	_ resource.ResourceWithImportState = &policyEventMatcherResource{}
)

func newPolicyEventMatcherResource() resource.Resource {
	return &policyEventMatcherResource{}
}

type policyEventMatcherResource struct {
	resourceBase
}

type policyEventMatcherModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ExecutionLogging types.Bool   `tfsdk:"execution_logging"`
	Action           types.String `tfsdk:"action"`
	ClientIP         types.String `tfsdk:"client_ip"`
	App              types.String `tfsdk:"app"`
	Model            types.String `tfsdk:"model"`
	Query            types.String `tfsdk:"query"`
}

func (r *policyEventMatcherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_event_matcher"
}

func (r *policyEventMatcherResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			// action is api.EventActions on the wire but SDKv2 declared no enum validator or
			// description for it, and the generated docs show it as a bare String. Kept that
			// way deliberately rather than "improving" it, which would change the doc page.
			"action": schema.StringAttribute{
				Optional: true,
			},
			"client_ip": schema.StringAttribute{
				Optional: true,
			},
			"app": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedAppEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedAppEnumEnumValues),
				},
			},
			"model": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedModelEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedModelEnumEnumValues),
				},
			},
			"query": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// toRequest leaves each Nullable wrapper untouched when the config value is null. All five
// of these fields are gated on IsSet() rather than nil-ness in the generated ToMap, so an
// untouched wrapper is what omits the key - matching SDKv2's `if v != "" { r.X.Set(...) }`
// guards. Calling Set(nil) instead would send an explicit JSON null, which is a different
// request. Same shape as authentik_stage_authenticator_webauthn's
// authenticator_attachment in batch 1.
func (r *policyEventMatcherResource) toRequest(data *policyEventMatcherModel) *api.EventMatcherPolicyRequest {
	body := &api.EventMatcherPolicyRequest{
		Name:             data.Name.ValueString(),
		ExecutionLogging: new(data.ExecutionLogging.ValueBool()),
	}

	if !data.Action.IsNull() {
		body.Action.Set(api.EventActions(data.Action.ValueString()).Ptr())
	}
	if !data.ClientIP.IsNull() {
		body.ClientIp.Set(data.ClientIP.ValueStringPointer())
	}
	if !data.App.IsNull() {
		body.App.Set(api.AppEnum(data.App.ValueString()).Ptr())
	}
	if !data.Model.IsNull() {
		body.Model.Set(api.ModelEnum(data.Model.ValueString()).Ptr())
	}
	if !data.Query.IsNull() {
		body.Query.Set(data.Query.ValueStringPointer())
	}
	return body
}

// fromAPI only assigns a field when the response actually carries it, mirroring SDKv2's
// `if res.HasX()` guards: an absent field leaves the value seeded from plan or prior state
// rather than nulling it. Where the field *is* present but null, null is written - which is
// the H1-correct answer for a genuinely nullable API field, and equivalent to what SDKv2
// did, since its legacy type system treated "" and null as the same thing.
func (r *policyEventMatcherResource) fromAPI(data *policyEventMatcherModel, res *api.EventMatcherPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())

	if res.HasAction() {
		if v := res.Action.Get(); v != nil {
			data.Action = types.StringValue(string(*v))
		} else {
			data.Action = types.StringNull()
		}
	}
	if res.HasClientIp() {
		data.ClientIP = helpers.StringPtrOrNull(res.ClientIp.Get())
	}
	if res.HasApp() {
		if v := res.App.Get(); v != nil {
			data.App = types.StringValue(string(*v))
		} else {
			data.App = types.StringNull()
		}
	}
	if res.HasModel() {
		if v := res.Model.Get(); v != nil {
			data.Model = types.StringValue(string(*v))
		} else {
			data.Model = types.StringNull()
		}
	}
	if res.HasQuery() {
		data.Query = helpers.StringPtrOrNull(res.Query.Get())
	}
}

func (r *policyEventMatcherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyEventMatcherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesEventMatcherCreate(ctx).EventMatcherPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyEventMatcherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyEventMatcherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesEventMatcherRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyEventMatcherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyEventMatcherModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesEventMatcherUpdate(ctx, data.ID.ValueString()).EventMatcherPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyEventMatcherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyEventMatcherModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesEventMatcherDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
