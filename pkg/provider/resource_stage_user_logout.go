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
	_ resource.Resource                = &stageUserLogoutResource{}
	_ resource.ResourceWithConfigure   = &stageUserLogoutResource{}
	_ resource.ResourceWithImportState = &stageUserLogoutResource{}
)

func newStageUserLogoutResource() resource.Resource {
	return &stageUserLogoutResource{}
}

type stageUserLogoutResource struct {
	resourceBase
}

type stageUserLogoutModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *stageUserLogoutResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_user_logout"
}

func (r *stageUserLogoutResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
		},
	}
}

func (r *stageUserLogoutResource) toRequest(data *stageUserLogoutModel) *api.UserLogoutStageRequest {
	return &api.UserLogoutStageRequest{
		Name: data.Name.ValueString(),
	}
}

func (r *stageUserLogoutResource) fromAPI(data *stageUserLogoutModel, res *api.UserLogoutStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
}

func (r *stageUserLogoutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageUserLogoutModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLogoutCreate(ctx).UserLogoutStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserLogoutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageUserLogoutModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLogoutRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageUserLogoutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageUserLogoutModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserLogoutUpdate(ctx, data.ID.ValueString()).UserLogoutStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserLogoutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageUserLogoutModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesUserLogoutDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
