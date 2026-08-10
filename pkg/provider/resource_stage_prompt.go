package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	_ resource.Resource                = &stagePromptResource{}
	_ resource.ResourceWithConfigure   = &stagePromptResource{}
	_ resource.ResourceWithImportState = &stagePromptResource{}
)

func newStagePromptResource() resource.Resource {
	return &stagePromptResource{}
}

type stagePromptResource struct {
	resourceBase
}

type stagePromptModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Fields             types.List   `tfsdk:"fields"`
	ValidationPolicies types.List   `tfsdk:"validation_policies"`
}

func (r *stagePromptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_prompt"
}

func (r *stagePromptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"fields": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"validation_policies": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *stagePromptResource) toRequest(ctx context.Context, data *stagePromptModel) (*api.PromptStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	fields, d := helpers.SliceOrEmpty[string](ctx, data.Fields)
	diags.Append(d...)

	// validation_policies is Optional, so this is the one that matters here: as a nil
	// slice it would be omitted and removing it from config would not clear it.
	validationPolicies, d := helpers.SliceOrEmpty[string](ctx, data.ValidationPolicies)
	diags.Append(d...)

	return &api.PromptStageRequest{
		Name:               data.Name.ValueString(),
		Fields:             fields,
		ValidationPolicies: validationPolicies,
	}, diags
}

func (r *stagePromptResource) fromAPI(ctx context.Context, data *stagePromptModel, res *api.PromptStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	fields, d := helpers.MergeStringList(ctx, data.Fields, res.Fields)
	diags.Append(d...)
	data.Fields = fields

	validationPolicies, d := helpers.MergeStringList(ctx, data.ValidationPolicies, res.ValidationPolicies)
	diags.Append(d...)
	data.ValidationPolicies = validationPolicies

	return diags
}

func (r *stagePromptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stagePromptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptStagesCreate(ctx).PromptStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePromptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stagePromptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptStagesRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePromptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stagePromptModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptStagesUpdate(ctx, data.ID.ValueString()).PromptStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePromptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stagePromptModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesPromptStagesDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
