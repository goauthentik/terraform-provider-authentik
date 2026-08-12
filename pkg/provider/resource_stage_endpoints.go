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
	_ resource.Resource                = &stageEndpointsResource{}
	_ resource.ResourceWithConfigure   = &stageEndpointsResource{}
	_ resource.ResourceWithImportState = &stageEndpointsResource{}
)

func newStageEndpointsResource() resource.Resource {
	return &stageEndpointsResource{}
}

type stageEndpointsResource struct {
	resourceBase
}

type stageEndpointsModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Connector types.String `tfsdk:"connector"`
	Mode      types.String `tfsdk:"mode"`
}

func (r *stageEndpointsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_endpoints"
}

func (r *stageEndpointsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	modeDefault := helpers.StringDefault(string(api.STAGEMODEENUM_OPTIONAL))

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
			"connector": schema.StringAttribute{
				Required: true,
			},
			"mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             modeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedStageModeEnumEnumValues), helpers.WithDefault(modeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedStageModeEnumEnumValues),
				},
			},
		},
	}
}

func (r *stageEndpointsResource) toRequest(data *stageEndpointsModel) *api.EndpointStageRequest {
	return &api.EndpointStageRequest{
		Name:      data.Name.ValueString(),
		Connector: data.Connector.ValueString(),
		Mode:      api.StageModeEnum(data.Mode.ValueString()).Ptr(),
	}
}

func (r *stageEndpointsResource) fromAPI(data *stageEndpointsModel, res *api.EndpointStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Connector = types.StringValue(res.Connector)
	if res.Mode != nil {
		data.Mode = types.StringValue(string(*res.Mode))
	}
}

func (r *stageEndpointsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageEndpointsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEndpointsCreate(ctx).EndpointStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageEndpointsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageEndpointsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEndpointsRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageEndpointsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageEndpointsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesEndpointsUpdate(ctx, data.ID.ValueString()).EndpointStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageEndpointsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageEndpointsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesEndpointsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
