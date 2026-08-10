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
	_ resource.Resource                = &stageRedirectResource{}
	_ resource.ResourceWithConfigure   = &stageRedirectResource{}
	_ resource.ResourceWithImportState = &stageRedirectResource{}
)

func newStageRedirectResource() resource.Resource {
	return &stageRedirectResource{}
}

type stageRedirectResource struct {
	resourceBase
}

type stageRedirectModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Mode         types.String `tfsdk:"mode"`
	KeepContext  types.Bool   `tfsdk:"keep_context"`
	TargetStatic types.String `tfsdk:"target_static"`
	TargetFlow   types.String `tfsdk:"target_flow"`
}

func (r *stageRedirectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_redirect"
}

func (r *stageRedirectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	modeDefault := helpers.StringDefault(string(api.REDIRECTSTAGEMODEENUM_FLOW))
	keepContextDefault := helpers.BoolDefault(true)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Flows & Stages --- ",
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
			"mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             modeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedRedirectStageModeEnumEnumValues), helpers.WithDefault(modeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedRedirectStageModeEnumEnumValues),
				},
			},
			"keep_context": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             keepContextDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(keepContextDefault.Value())),
			},
			"target_static": schema.StringAttribute{
				Optional: true,
			},
			"target_flow": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *stageRedirectResource) toRequest(data *stageRedirectModel) *api.RedirectStageRequest {
	return &api.RedirectStageRequest{
		Name:         data.Name.ValueString(),
		Mode:         api.RedirectStageModeEnum(data.Mode.ValueString()),
		KeepContext:  new(data.KeepContext.ValueBool()),
		TargetStatic: helpers.StringPtr(data.TargetStatic),
		TargetFlow:   *api.NewNullableString(helpers.StringPtr(data.TargetFlow)),
	}
}

// fromAPI splits its mapping by whether the attribute has a schema Default. An
// attribute with a Default can never legitimately be null in state, so it takes the
// API value verbatim; the prior-aware helpers are only for attributes that can be
// null. Using BoolOrNull on keep_context would break `terraform import`, where prior
// state is null: BoolOrNull(null, false) returns null, and the next plan would then
// apply the default and report a spurious null -> false update.
func (r *stageRedirectResource) fromAPI(data *stageRedirectModel, res *api.RedirectStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Mode = types.StringValue(string(res.Mode))
	data.KeepContext = types.BoolValue(res.GetKeepContext())
	data.TargetStatic = helpers.StringOrNull(data.TargetStatic, res.GetTargetStatic())
	data.TargetFlow = helpers.StringPtrOrNull(res.TargetFlow.Get())
}

func (r *stageRedirectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageRedirectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesRedirectCreate(ctx).RedirectStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageRedirectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageRedirectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesRedirectRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageRedirectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageRedirectModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesRedirectUpdate(ctx, data.ID.ValueString()).RedirectStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageRedirectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageRedirectModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesRedirectDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
