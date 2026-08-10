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
	_ resource.Resource                = &stageAccountLockdownResource{}
	_ resource.ResourceWithConfigure   = &stageAccountLockdownResource{}
	_ resource.ResourceWithImportState = &stageAccountLockdownResource{}
)

func newStageAccountLockdownResource() resource.Resource {
	return &stageAccountLockdownResource{}
}

type stageAccountLockdownResource struct {
	resourceBase
}

type stageAccountLockdownModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	DeactivateUser            types.Bool   `tfsdk:"deactivate_user"`
	SetUnusablePassword       types.Bool   `tfsdk:"set_unusable_password"`
	DeleteSessions            types.Bool   `tfsdk:"delete_sessions"`
	RevokeTokens              types.Bool   `tfsdk:"revoke_tokens"`
	SelfServiceCompletionFlow types.String `tfsdk:"self_service_completion_flow"`
}

func (r *stageAccountLockdownResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_account_lockdown"
}

func (r *stageAccountLockdownResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	deactivateUserDefault := helpers.BoolDefault(true)
	setUnusablePasswordDefault := helpers.BoolDefault(true)
	deleteSessionsDefault := helpers.BoolDefault(true)
	revokeTokensDefault := helpers.BoolDefault(true)

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
			"deactivate_user": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             deactivateUserDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(deactivateUserDefault.Value())),
			},
			"set_unusable_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             setUnusablePasswordDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(setUnusablePasswordDefault.Value())),
			},
			"delete_sessions": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             deleteSessionsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(deleteSessionsDefault.Value())),
			},
			"revoke_tokens": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             revokeTokensDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(revokeTokensDefault.Value())),
			},
			"self_service_completion_flow": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *stageAccountLockdownResource) toRequest(data *stageAccountLockdownModel) *api.AccountLockdownStageRequest {
	return &api.AccountLockdownStageRequest{
		Name:                      data.Name.ValueString(),
		DeactivateUser:            new(data.DeactivateUser.ValueBool()),
		SetUnusablePassword:       new(data.SetUnusablePassword.ValueBool()),
		DeleteSessions:            new(data.DeleteSessions.ValueBool()),
		RevokeTokens:              new(data.RevokeTokens.ValueBool()),
		SelfServiceCompletionFlow: *api.NewNullableString(helpers.StringPtr(data.SelfServiceCompletionFlow)),
	}
}

func (r *stageAccountLockdownResource) fromAPI(data *stageAccountLockdownModel, res *api.AccountLockdownStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	// All four bools have a schema Default, so they can never legitimately be null in
	// state and must take the API value verbatim - see the comment on
	// stageRedirectResource.fromAPI for why the prior-aware helpers break import here.
	data.DeactivateUser = types.BoolValue(res.GetDeactivateUser())
	data.SetUnusablePassword = types.BoolValue(res.GetSetUnusablePassword())
	data.DeleteSessions = types.BoolValue(res.GetDeleteSessions())
	data.RevokeTokens = types.BoolValue(res.GetRevokeTokens())
	data.SelfServiceCompletionFlow = helpers.StringPtrOrNull(res.SelfServiceCompletionFlow.Get())
}

func (r *stageAccountLockdownResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageAccountLockdownModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAccountLockdownCreate(ctx).AccountLockdownStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAccountLockdownResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageAccountLockdownModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAccountLockdownRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageAccountLockdownResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageAccountLockdownModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesAccountLockdownUpdate(ctx, data.ID.ValueString()).AccountLockdownStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageAccountLockdownResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageAccountLockdownModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesAccountLockdownDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
