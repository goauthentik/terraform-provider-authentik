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
	_ resource.Resource                = &stagePromptFieldResource{}
	_ resource.ResourceWithConfigure   = &stagePromptFieldResource{}
	_ resource.ResourceWithImportState = &stagePromptFieldResource{}
)

func newStagePromptFieldResource() resource.Resource {
	return &stagePromptFieldResource{}
}

type stagePromptFieldResource struct {
	resourceBase
}

type stagePromptFieldModel struct {
	ID                     types.String            `tfsdk:"id"`
	Name                   types.String            `tfsdk:"name"`
	FieldKey               types.String            `tfsdk:"field_key"`
	Label                  types.String            `tfsdk:"label"`
	Type                   types.String            `tfsdk:"type"`
	Required               types.Bool              `tfsdk:"required"`
	Placeholder            helpers.ExpressionValue `tfsdk:"placeholder"`
	PlaceholderExpression  types.Bool              `tfsdk:"placeholder_expression"`
	InitialValue           helpers.ExpressionValue `tfsdk:"initial_value"`
	InitialValueExpression types.Bool              `tfsdk:"initial_value_expression"`
	Order                  types.Int32             `tfsdk:"order"`
	SubText                types.String            `tfsdk:"sub_text"`
}

func (r *stagePromptFieldResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_prompt_field"
}

func (r *stagePromptFieldResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiredDefault := helpers.BoolDefault(false)
	placeholderExpressionDefault := helpers.BoolDefault(false)
	initialValueExpressionDefault := helpers.BoolDefault(false)
	subTextDefault := helpers.StringDefault("")

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
			"field_key": schema.StringAttribute{
				Required: true,
			},
			"label": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: helpers.EnumToDescription(api.AllowedPromptTypeEnumEnumValues),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedPromptTypeEnumEnumValues),
				},
			},
			"required": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             requiredDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(requiredDefault.Value())),
			},
			// placeholder and initial_value carried DiffSuppressExpression in SDKv2 and are
			// the first two of the 18 such attributes to migrate. helpers.ExpressionType
			// replaces the DiffSuppressFunc with real semantic equality: these fields can
			// hold an authentik expression, authentik returns expressions with trailing
			// newlines stripped, and a heredoc config value ending in "\n" would otherwise
			// show a permanent diff. Both stay protocol `string`, so the doc page and the
			// state layout are unchanged.
			"placeholder": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Optional:   true,
			},
			"placeholder_expression": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             placeholderExpressionDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(placeholderExpressionDefault.Value())),
			},
			"initial_value": schema.StringAttribute{
				CustomType: helpers.ExpressionType{},
				Optional:   true,
			},
			"initial_value_expression": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             initialValueExpressionDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(initialValueExpressionDefault.Value())),
			},
			"order": schema.Int32Attribute{
				Optional: true,
			},
			"sub_text": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             subTextDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(subTextDefault.Value())),
			},
		},
	}
}

func (r *stagePromptFieldResource) toRequest(data *stagePromptFieldModel) *api.PromptRequest {
	return &api.PromptRequest{
		Name:                   data.Name.ValueString(),
		FieldKey:               data.FieldKey.ValueString(),
		Label:                  data.Label.ValueString(),
		Type:                   api.PromptTypeEnum(data.Type.ValueString()),
		Required:               new(data.Required.ValueBool()),
		PlaceholderExpression:  new(data.PlaceholderExpression.ValueBool()),
		InitialValueExpression: new(data.InitialValueExpression.ValueBool()),
		SubText:                new(data.SubText.ValueString()),
		Placeholder:            helpers.ExpressionPtr(data.Placeholder),
		InitialValue:           helpers.ExpressionPtr(data.InitialValue),
		Order:                  helpers.Int32Ptr(data.Order),
	}
}

func (r *stagePromptFieldResource) fromAPI(data *stagePromptFieldModel, res *api.Prompt) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.FieldKey = types.StringValue(res.FieldKey)
	data.Label = types.StringValue(res.Label)
	data.Type = types.StringValue(string(res.Type))
	// No Defaults on these three, so they keep the prior-aware helpers. Semantic equality
	// deliberately does not cover the null-vs-"" case: value_semantic_equality.go returns
	// early when the prior value is null, so ExpressionOrNull is still required on top of
	// the custom type.
	data.Placeholder = helpers.ExpressionOrNull(data.Placeholder, res.GetPlaceholder())
	data.InitialValue = helpers.ExpressionOrNull(data.InitialValue, res.GetInitialValue())
	data.Order = helpers.Int32OrNull(data.Order, res.GetOrder())
	// These four have Defaults and take the API value verbatim.
	data.Required = types.BoolValue(res.GetRequired())
	data.PlaceholderExpression = types.BoolValue(res.GetPlaceholderExpression())
	data.InitialValueExpression = types.BoolValue(res.GetInitialValueExpression())
	data.SubText = types.StringValue(res.GetSubText())
}

func (r *stagePromptFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stagePromptFieldModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptPromptsCreate(ctx).PromptRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePromptFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stagePromptFieldModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptPromptsRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stagePromptFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stagePromptFieldModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPromptPromptsUpdate(ctx, data.ID.ValueString()).PromptRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePromptFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stagePromptFieldModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesPromptPromptsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
