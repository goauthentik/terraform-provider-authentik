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
	_ resource.Resource                = &stageConsentResource{}
	_ resource.ResourceWithConfigure   = &stageConsentResource{}
	_ resource.ResourceWithImportState = &stageConsentResource{}
)

func newStageConsentResource() resource.Resource {
	return &stageConsentResource{}
}

type stageConsentResource struct {
	resourceBase
}

type stageConsentModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Mode            types.String `tfsdk:"mode"`
	ConsentExpireIn types.String `tfsdk:"consent_expire_in"`
}

func (r *stageConsentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_consent"
}

func (r *stageConsentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	modeDefault := helpers.StringDefault(string(api.CONSENTMODEENUM_ALWAYS_REQUIRE))
	expireInDefault := helpers.StringDefault("weeks=4")

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
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedConsentModeEnumEnumValues), helpers.WithDefault(modeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedConsentModeEnumEnumValues),
				},
			},
			"consent_expire_in": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             expireInDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(expireInDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
		},
	}
}

func (r *stageConsentResource) toRequest(data *stageConsentModel) *api.ConsentStageRequest {
	req := &api.ConsentStageRequest{
		Name:            data.Name.ValueString(),
		ConsentExpireIn: helpers.StringPtr(data.ConsentExpireIn),
	}
	if !data.Mode.IsNull() {
		req.Mode = api.ConsentModeEnum(data.Mode.ValueString()).Ptr()
	}
	return req
}

func (r *stageConsentResource) fromAPI(data *stageConsentModel, res *api.ConsentStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	if res.Mode != nil {
		data.Mode = types.StringValue(string(*res.Mode))
	}
	data.ConsentExpireIn = types.StringValue(res.GetConsentExpireIn())
}

func (r *stageConsentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageConsentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesConsentCreate(ctx).ConsentStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageConsentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageConsentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesConsentRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageConsentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageConsentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesConsentUpdate(ctx, data.ID.ValueString()).ConsentStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageConsentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageConsentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesConsentDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
