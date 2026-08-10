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
	_ resource.Resource                = &stageAuthenticatorStaticResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorStaticResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorStaticResource{}
)

func newStageAuthenticatorStaticResource() resource.Resource {
	return &stageAuthenticatorStaticResource{}
}

type stageAuthenticatorStaticResource struct {
	resourceBase
}

type stageAuthenticatorStaticModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	FriendlyName  types.String `tfsdk:"friendly_name"`
	ConfigureFlow types.String `tfsdk:"configure_flow"`
	TokenCount    types.Int32  `tfsdk:"token_count"`
	TokenLength   types.Int32  `tfsdk:"token_length"`
}

func (r *stageAuthenticatorStaticResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_static"
}

func (r *stageAuthenticatorStaticResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")
	tokenCountDefault := helpers.Int32Default(6)
	tokenLengthDefault := helpers.Int32Default(12)

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
			"friendly_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             friendlyNameDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(friendlyNameDefault.Value())),
			},
			"configure_flow": schema.StringAttribute{
				Optional: true,
			},
			"token_count": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             tokenCountDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(tokenCountDefault.Value())),
			},
			"token_length": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             tokenLengthDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(tokenLengthDefault.Value())),
			},
		},
	}
}

func (r *stageAuthenticatorStaticResource) toRequest(data *stageAuthenticatorStaticModel) *api.AuthenticatorStaticStageRequest {
	return &api.AuthenticatorStaticStageRequest{
		Name:          data.Name.ValueString(),
		TokenCount:    new(data.TokenCount.ValueInt32()),
		TokenLength:   new(data.TokenLength.ValueInt32()),
		FriendlyName:  helpers.StringPtr(data.FriendlyName),
		ConfigureFlow: *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
	}
}

func (r *stageAuthenticatorStaticResource) fromAPI(data *stageAuthenticatorStaticModel, res *api.AuthenticatorStaticStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.TokenCount = types.Int32Value(res.GetTokenCount())
	data.TokenLength = types.Int32Value(res.GetTokenLength())
	data.FriendlyName = types.StringValue(res.GetFriendlyName())
	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
}

func (r *stageAuthenticatorStaticResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorStaticModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorStaticCreate(ctx).AuthenticatorStaticStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorStaticResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorStaticModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorStaticRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorStaticResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorStaticModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorStaticUpdate(ctx, data.ID.ValueString()).AuthenticatorStaticStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorStaticResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorStaticModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorStaticDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
