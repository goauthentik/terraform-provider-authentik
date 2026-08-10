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
	_ resource.Resource                = &stageDenyResource{}
	_ resource.ResourceWithConfigure   = &stageDenyResource{}
	_ resource.ResourceWithImportState = &stageDenyResource{}
)

func newStageDenyResource() resource.Resource {
	return &stageDenyResource{}
}

type stageDenyResource struct {
	resourceBase
}

type stageDenyModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DenyMessage types.String `tfsdk:"deny_message"`
}

func (r *stageDenyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_deny"
}

func (r *stageDenyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"deny_message": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *stageDenyResource) toRequest(data *stageDenyModel) *api.DenyStageRequest {
	return &api.DenyStageRequest{
		Name:        data.Name.ValueString(),
		DenyMessage: helpers.StringPtr(data.DenyMessage),
	}
}

func (r *stageDenyResource) fromAPI(data *stageDenyModel, res *api.DenyStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.DenyMessage = helpers.StringOrNull(data.DenyMessage, res.GetDenyMessage())
}

func (r *stageDenyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageDenyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesDenyCreate(ctx).DenyStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageDenyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageDenyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesDenyRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageDenyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageDenyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesDenyUpdate(ctx, data.ID.ValueString()).DenyStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageDenyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageDenyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesDenyDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
