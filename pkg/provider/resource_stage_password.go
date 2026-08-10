package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = &stagePasswordResource{}
	_ resource.ResourceWithConfigure   = &stagePasswordResource{}
	_ resource.ResourceWithImportState = &stagePasswordResource{}
)

func newStagePasswordResource() resource.Resource {
	return &stagePasswordResource{}
}

type stagePasswordResource struct {
	resourceBase
}

type stagePasswordModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Backends                   types.List   `tfsdk:"backends"`
	ConfigureFlow              types.String `tfsdk:"configure_flow"`
	FailedAttemptsBeforeCancel types.Int32  `tfsdk:"failed_attempts_before_cancel"`
	AllowShowPassword          types.Bool   `tfsdk:"allow_show_password"`
}

func (r *stagePasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_password"
}

func (r *stagePasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	failedAttemptsDefault := helpers.Int32Default(5)
	allowShowPasswordDefault := helpers.BoolDefault(false)

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
			// SDKv2 declared the enum description and validator on the list's Elem
			// schema, but CoreConfigSchema drops an element schema's description for
			// primitive lists - docs/resources/stage_password.md renders `backends`
			// with no description at all. Keep the element-level validation (as
			// listvalidator.ValueStringsAre) but deliberately no MarkdownDescription,
			// or the docs-drift gate fires.
			"backends": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(helpers.OneOf(api.AllowedBackendsEnumEnumValues)),
				},
			},
			"configure_flow": schema.StringAttribute{
				Optional: true,
			},
			"failed_attempts_before_cancel": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             failedAttemptsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(failedAttemptsDefault.Value())),
			},
			"allow_show_password": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             allowShowPasswordDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(allowShowPasswordDefault.Value())),
			},
		},
	}
}

func (r *stagePasswordResource) toRequest(ctx context.Context, data *stagePasswordModel) (*api.PasswordStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// backends is Required, so it can never be null - SliceOrEmpty is used anyway so
	// every list in this package goes through one code path.
	backends, d := helpers.SliceOrEmpty[string](ctx, data.Backends)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.PasswordStageRequest{
		Name:                       data.Name.ValueString(),
		Backends:                   helpers.CastSliceString[api.BackendsEnum](backends),
		AllowShowPassword:          new(data.AllowShowPassword.ValueBool()),
		ConfigureFlow:              *api.NewNullableString(helpers.StringPtr(data.ConfigureFlow)),
		FailedAttemptsBeforeCancel: new(data.FailedAttemptsBeforeCancel.ValueInt32()),
	}, diags
}

// fromAPI merges backends rather than overwriting it. SDKv2 wrote the API's order
// straight into state, which was tolerable because core demoted post-apply
// inconsistencies to warnings for legacy SDK providers; the framework makes them hard
// errors, and `backends` is Required, so any reordering by the API would fail the
// apply. MergeStringList keeps the configured order and is a no-op when the orders
// already agree.
func (r *stagePasswordResource) fromAPI(ctx context.Context, data *stagePasswordModel, res *api.PasswordStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)

	backends := make([]string, len(res.Backends))
	for i, b := range res.Backends {
		backends[i] = string(b)
	}
	merged, d := helpers.MergeStringList(ctx, data.Backends, backends)
	diags.Append(d...)
	data.Backends = merged

	data.ConfigureFlow = helpers.StringPtrOrNull(res.ConfigureFlow.Get())
	// Both of these have a schema Default, so they can never legitimately be null in
	// state and must take the API value verbatim - see the comment on
	// stageRedirectResource.fromAPI for why the prior-aware helpers break import here.
	data.FailedAttemptsBeforeCancel = types.Int32Value(res.GetFailedAttemptsBeforeCancel())
	data.AllowShowPassword = types.BoolValue(res.GetAllowShowPassword())

	return diags
}

func (r *stagePasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stagePasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPasswordCreate(ctx).PasswordStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stagePasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPasswordRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stagePasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesPasswordUpdate(ctx, data.ID.ValueString()).PasswordStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stagePasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stagePasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesPasswordDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
