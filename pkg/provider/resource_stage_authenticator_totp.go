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
	_ resource.Resource                = &stageAuthenticatorTOTPResource{}
	_ resource.ResourceWithConfigure   = &stageAuthenticatorTOTPResource{}
	_ resource.ResourceWithImportState = &stageAuthenticatorTOTPResource{}
)

func newStageAuthenticatorTOTPResource() resource.Resource {
	return &stageAuthenticatorTOTPResource{}
}

type stageAuthenticatorTOTPResource struct {
	resourceBase
}

type stageAuthenticatorTOTPModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	FriendlyName  types.String `tfsdk:"friendly_name"`
	ConfigureFlow types.String `tfsdk:"configure_flow"`
	Digits        types.String `tfsdk:"digits"`
}

func (r *stageAuthenticatorTOTPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_authenticator_totp"
}

func (r *stageAuthenticatorTOTPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	friendlyNameDefault := helpers.StringDefault("")
	digitsDefault := helpers.StringDefault(string(api.DIGITSENUM__6))

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
			"digits": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             digitsDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedDigitsEnumEnumValues), helpers.WithDefault(digitsDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedDigitsEnumEnumValues),
				},
			},
		},
	}
}

func (r *stageAuthenticatorTOTPResource) toRequest(data *stageAuthenticatorTOTPModel) *api.AuthenticatorTOTPStageRequest {
	return &api.AuthenticatorTOTPStageRequest{
		Name:          data.Name.ValueString(),
		Digits:        api.DigitsEnum(data.Digits.ValueString()),
		FriendlyName:  helpers.StringPtr(data.FriendlyName),
		ConfigureFlow: *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
	}
}

func (r *stageAuthenticatorTOTPResource) fromAPI(data *stageAuthenticatorTOTPModel, res *api.AuthenticatorTOTPStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Digits = types.StringValue(string(res.Digits))
	data.FriendlyName = helpers.StringOrNull(data.FriendlyName, res.GetFriendlyName())
	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
}

func (r *stageAuthenticatorTOTPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAuthenticatorTOTPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorTotpCreate(ctx).AuthenticatorTOTPStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorTOTPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAuthenticatorTOTPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorTotpRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAuthenticatorTOTPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAuthenticatorTOTPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAuthenticatorTotpUpdate(ctx, data.ID.ValueString()).AuthenticatorTOTPStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAuthenticatorTOTPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAuthenticatorTOTPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAuthenticatorTotpDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
