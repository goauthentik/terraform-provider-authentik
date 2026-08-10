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
	_ resource.Resource                = &policyUniquePasswordResource{}
	_ resource.ResourceWithConfigure   = &policyUniquePasswordResource{}
	_ resource.ResourceWithImportState = &policyUniquePasswordResource{}
)

func newPolicyUniquePasswordResource() resource.Resource {
	return &policyUniquePasswordResource{}
}

type policyUniquePasswordResource struct {
	resourceBase
}

type policyUniquePasswordModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ExecutionLogging       types.Bool   `tfsdk:"execution_logging"`
	PasswordField          types.String `tfsdk:"password_field"`
	NumHistoricalPasswords types.Int32  `tfsdk:"num_historical_passwords"`
}

func (r *policyUniquePasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_unique_password"
}

func (r *policyUniquePasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	executionLoggingDefault := helpers.BoolDefault(false)
	passwordFieldDefault := helpers.StringDefault("password")
	numHistoricalPasswordsDefault := helpers.Int32Default(1)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Customization --- ",
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
			"execution_logging": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             executionLoggingDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(executionLoggingDefault.Value())),
			},
			"password_field": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             passwordFieldDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(passwordFieldDefault.Value())),
			},
			"num_historical_passwords": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             numHistoricalPasswordsDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(numHistoricalPasswordsDefault.Value())),
			},
		},
	}
}

func (r *policyUniquePasswordResource) toRequest(data *policyUniquePasswordModel) *api.UniquePasswordPolicyRequest {
	return &api.UniquePasswordPolicyRequest{
		Name:                   data.Name.ValueString(),
		ExecutionLogging:       new(data.ExecutionLogging.ValueBool()),
		PasswordField:          new(data.PasswordField.ValueString()),
		NumHistoricalPasswords: new(data.NumHistoricalPasswords.ValueInt32()),
	}
}

func (r *policyUniquePasswordResource) fromAPI(data *policyUniquePasswordModel, res *api.UniquePasswordPolicy) {
	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.ExecutionLogging = types.BoolValue(res.GetExecutionLogging())
	data.PasswordField = types.StringValue(res.GetPasswordField())
	data.NumHistoricalPasswords = types.Int32Value(res.GetNumHistoricalPasswords())
}

func (r *policyUniquePasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data policyUniquePasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesUniquePasswordCreate(ctx).UniquePasswordPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyUniquePasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data policyUniquePasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesUniquePasswordRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *policyUniquePasswordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data policyUniquePasswordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.PoliciesAPI.PoliciesUniquePasswordUpdate(ctx, data.ID.ValueString()).UniquePasswordPolicyRequest(*r.toRequest(&data)).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	r.fromAPI(&data, res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *policyUniquePasswordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data policyUniquePasswordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.PoliciesAPI.PoliciesUniquePasswordDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
