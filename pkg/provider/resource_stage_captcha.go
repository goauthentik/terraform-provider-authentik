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
	_ resource.Resource                = &stageCaptchaResource{}
	_ resource.ResourceWithConfigure   = &stageCaptchaResource{}
	_ resource.ResourceWithImportState = &stageCaptchaResource{}
)

func newStageCaptchaResource() resource.Resource {
	return &stageCaptchaResource{}
}

type stageCaptchaResource struct {
	resourceBase
}

type stageCaptchaModel struct {
	ID                  types.String  `tfsdk:"id"`
	Name                types.String  `tfsdk:"name"`
	PublicKey           types.String  `tfsdk:"public_key"`
	PrivateKey          types.String  `tfsdk:"private_key"`
	JsURL               types.String  `tfsdk:"js_url"`
	APIURL              types.String  `tfsdk:"api_url"`
	ScoreMinThreshold   types.Float64 `tfsdk:"score_min_threshold"`
	ScoreMaxThreshold   types.Float64 `tfsdk:"score_max_threshold"`
	ErrorOnInvalidScore types.Bool    `tfsdk:"error_on_invalid_score"`
	Interactive         types.Bool    `tfsdk:"interactive"`
}

func (r *stageCaptchaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_captcha"
}

func (r *stageCaptchaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsURLDefault := helpers.StringDefault("https://www.recaptcha.net/recaptcha/api.js")
	apiURLDefault := helpers.StringDefault("https://www.recaptcha.net/recaptcha/api/siteverify")
	// The only two TypeFloat attributes in the whole provider, alongside
	// authentik_stage_authenticator_validate's. Both TypeFloat and Float64Attribute
	// serialise to protocol `number`, so this is a free change. Desc renders these with
	// %v, which prints 1 rather than 1.0 - matching the existing doc page exactly.
	scoreMinDefault := helpers.Float64Default(0.5)
	scoreMaxDefault := helpers.Float64Default(1)
	errorOnInvalidScoreDefault := helpers.BoolDefault(true)
	interactiveDefault := helpers.BoolDefault(false)

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
			"public_key": schema.StringAttribute{
				Required: true,
			},
			"private_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"js_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             jsURLDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(jsURLDefault.Value())),
			},
			"api_url": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             apiURLDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(apiURLDefault.Value())),
			},
			"score_min_threshold": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             scoreMinDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(scoreMinDefault.Value())),
			},
			"score_max_threshold": schema.Float64Attribute{
				Optional:            true,
				Computed:            true,
				Default:             scoreMaxDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(scoreMaxDefault.Value())),
			},
			"error_on_invalid_score": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             errorOnInvalidScoreDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(errorOnInvalidScoreDefault.Value())),
			},
			"interactive": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             interactiveDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(interactiveDefault.Value())),
			},
		},
	}
}

func (r *stageCaptchaResource) toRequest(data *stageCaptchaModel) *api.CaptchaStageRequest {
	return &api.CaptchaStageRequest{
		Name:                data.Name.ValueString(),
		PublicKey:           data.PublicKey.ValueString(),
		PrivateKey:          data.PrivateKey.ValueString(),
		JsUrl:               helpers.StringPtr(data.JsURL),
		ApiUrl:              helpers.StringPtr(data.APIURL),
		ScoreMinThreshold:   new(data.ScoreMinThreshold.ValueFloat64()),
		ScoreMaxThreshold:   new(data.ScoreMaxThreshold.ValueFloat64()),
		ErrorOnInvalidScore: new(data.ErrorOnInvalidScore.ValueBool()),
		Interactive:         new(data.Interactive.ValueBool()),
	}
}

// fromAPI deliberately never assigns PrivateKey. This is one of the 18 H4 resources: the
// API never returns private_key (it is write-only server-side), and SDKv2's Read simply
// did not call SetWrapper for it, leaving whatever was already in the ResourceData.
// resp.State.Set is wholesale rather than incremental, so the equivalent here is that
// every caller seeds data from req.Plan (Create/Update) or req.State (Read) *before*
// calling this - the omission then preserves the configured secret instead of nulling it.
// Assigning types.StringNull() here, or seeding from an empty model, would wipe a Required
// attribute out of state and fail the apply.
func (r *stageCaptchaResource) fromAPI(data *stageCaptchaModel, res *api.CaptchaStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.PublicKey = types.StringValue(res.PublicKey)
	// Everything below carries a schema Default, so it takes the API value verbatim
	// rather than going through the prior-aware helpers.
	data.JsURL = types.StringValue(res.GetJsUrl())
	data.APIURL = types.StringValue(res.GetApiUrl())
	data.ScoreMinThreshold = types.Float64Value(res.GetScoreMinThreshold())
	data.ScoreMaxThreshold = types.Float64Value(res.GetScoreMaxThreshold())
	data.ErrorOnInvalidScore = types.BoolValue(res.GetErrorOnInvalidScore())
	data.Interactive = types.BoolValue(res.GetInteractive())
}

func (r *stageCaptchaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageCaptchaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesCaptchaCreate(ctx).CaptchaStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageCaptchaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	// Seeding from prior state is what carries private_key across a Read - see fromAPI.
	var data stageCaptchaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesCaptchaRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageCaptchaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageCaptchaModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesCaptchaUpdate(ctx, data.ID.ValueString()).CaptchaStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageCaptchaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageCaptchaModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesCaptchaDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
