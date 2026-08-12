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
	_ resource.Resource                = &stageSourceResource{}
	_ resource.ResourceWithConfigure   = &stageSourceResource{}
	_ resource.ResourceWithImportState = &stageSourceResource{}
)

func newStageSourceResource() resource.Resource {
	return &stageSourceResource{}
}

type stageSourceResource struct {
	resourceBase
}

type stageSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Source        types.String `tfsdk:"source"`
	ResumeTimeout types.String `tfsdk:"resume_timeout"`
}

func (r *stageSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_source"
}

func (r *stageSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resumeTimeoutDefault := helpers.StringDefault("minutes=10")

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
			"source": schema.StringAttribute{
				Optional: true,
			},
			"resume_timeout": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             resumeTimeoutDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(resumeTimeoutDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
		},
	}
}

func (r *stageSourceResource) toRequest(data *stageSourceModel) *api.SourceStageRequest {
	return &api.SourceStageRequest{
		Name:          data.Name.ValueString(),
		Source:        data.Source.ValueString(),
		ResumeTimeout: new(data.ResumeTimeout.ValueString()),
	}
}

func (r *stageSourceResource) fromAPI(data *stageSourceModel, res *api.SourceStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Source = helpers.StringOrNull(data.Source, res.Source)
	data.ResumeTimeout = types.StringValue(res.GetResumeTimeout())
}

func (r *stageSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesSourceCreate(ctx).SourceStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesSourceRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesSourceUpdate(ctx, data.ID.ValueString()).SourceStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageSourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesSourceDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
