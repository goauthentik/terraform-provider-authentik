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
	_ resource.Resource                = &stageUserWriteResource{}
	_ resource.ResourceWithConfigure   = &stageUserWriteResource{}
	_ resource.ResourceWithImportState = &stageUserWriteResource{}
)

// userWriteUserTypes is the subset of api.UserTypeEnum this stage accepts. SDKv2 spelled
// the same list out twice (once for the description, once for the validator); keeping it
// in one place here is what makes the generated doc page byte-identical.
var userWriteUserTypes = []api.UserTypeEnum{
	api.USERTYPEENUM_INTERNAL,
	api.USERTYPEENUM_EXTERNAL,
	api.USERTYPEENUM_SERVICE_ACCOUNT,
}

func newStageUserWriteResource() resource.Resource {
	return &stageUserWriteResource{}
}

type stageUserWriteResource struct {
	resourceBase
}

type stageUserWriteModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	CreateUsersAsInactive types.Bool   `tfsdk:"create_users_as_inactive"`
	UserCreationMode      types.String `tfsdk:"user_creation_mode"`
	CreateUsersGroup      types.String `tfsdk:"create_users_group"`
	UserPathTemplate      types.String `tfsdk:"user_path_template"`
	UserType              types.String `tfsdk:"user_type"`
}

func (r *stageUserWriteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_user_write"
}

func (r *stageUserWriteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	createUsersAsInactiveDefault := helpers.BoolDefault(true)
	userCreationModeDefault := helpers.StringDefault(string(api.USERCREATIONMODEENUM_CREATE_WHEN_REQUIRED))
	userPathTemplateDefault := helpers.StringDefault("")
	userTypeDefault := helpers.StringDefault(string(api.USERTYPEENUM_EXTERNAL))

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
			"create_users_as_inactive": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             createUsersAsInactiveDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(createUsersAsInactiveDefault.Value())),
			},
			"user_creation_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userCreationModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserCreationModeEnumEnumValues), helpers.WithDefault(userCreationModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserCreationModeEnumEnumValues),
				},
			},
			"create_users_group": schema.StringAttribute{
				Optional: true,
			},
			"user_path_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userPathTemplateDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(userPathTemplateDefault.Value())),
			},
			"user_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userTypeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(userWriteUserTypes), helpers.WithDefault(userTypeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(userWriteUserTypes),
				},
			},
		},
	}
}

func (r *stageUserWriteResource) toRequest(data *stageUserWriteModel) *api.UserWriteStageRequest {
	return &api.UserWriteStageRequest{
		Name:                  data.Name.ValueString(),
		CreateUsersAsInactive: new(data.CreateUsersAsInactive.ValueBool()),
		UserPathTemplate:      new(data.UserPathTemplate.ValueString()),
		UserCreationMode:      api.UserCreationModeEnum(data.UserCreationMode.ValueString()).Ptr(),
		UserType:              api.UserTypeEnum(data.UserType.ValueString()).Ptr(),
		CreateUsersGroup:      *api.NewNullableString(helpers.StringPtr(data.CreateUsersGroup)),
	}
}

func (r *stageUserWriteResource) fromAPI(data *stageUserWriteModel, res *api.UserWriteStage) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.CreateUsersGroup = helpers.StringPtrOrNull(res.CreateUsersGroup.Get())
	// Everything below has a schema Default, so it takes the API value verbatim.
	// user_path_template is the case that makes this matter: its default is "", which is
	// also the value the API returns when it is unset, so the prior-aware helper would
	// import it as null and then plan a spurious null -> "" update.
	data.CreateUsersAsInactive = types.BoolValue(res.GetCreateUsersAsInactive())
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.UserCreationMode = types.StringValue(string(res.GetUserCreationMode()))
	data.UserType = types.StringValue(string(res.GetUserType()))
}

func (r *stageUserWriteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageUserWriteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserWriteCreate(ctx).UserWriteStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserWriteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageUserWriteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserWriteRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageUserWriteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageUserWriteModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesUserWriteUpdate(ctx, data.ID.ValueString()).UserWriteStageRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageUserWriteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageUserWriteModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesUserWriteDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
