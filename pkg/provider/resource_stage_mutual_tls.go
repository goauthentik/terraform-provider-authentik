package provider

import (
	"context"

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
	_ resource.Resource                = &stageMutualTLSResource{}
	_ resource.ResourceWithConfigure   = &stageMutualTLSResource{}
	_ resource.ResourceWithImportState = &stageMutualTLSResource{}
)

func newStageMutualTLSResource() resource.Resource {
	return &stageMutualTLSResource{}
}

type stageMutualTLSResource struct {
	resourceBase
}

type stageMutualTLSModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Mode                   types.String `tfsdk:"mode"`
	CertAttribute          types.String `tfsdk:"cert_attribute"`
	UserAttribute          types.String `tfsdk:"user_attribute"`
	CertificateAuthorities types.List   `tfsdk:"certificate_authorities"`
}

func (r *stageMutualTLSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stage_mutual_tls"
}

func (r *stageMutualTLSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	modeDefault := helpers.StringDefault(string(api.STAGEMODEENUM_OPTIONAL))
	certAttributeDefault := helpers.StringDefault(string(api.CERTATTRIBUTEENUM_EMAIL))
	userAttributeDefault := helpers.StringDefault(string(api.USERATTRIBUTEENUM_EMAIL))

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
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedStageModeEnumEnumValues), helpers.WithDefault(modeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedStageModeEnumEnumValues),
				},
			},
			"cert_attribute": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             certAttributeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedCertAttributeEnumEnumValues), helpers.WithDefault(certAttributeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedCertAttributeEnumEnumValues),
				},
			},
			"user_attribute": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userAttributeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserAttributeEnumEnumValues), helpers.WithDefault(userAttributeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserAttributeEnumEnumValues),
				},
			},
			"certificate_authorities": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *stageMutualTLSResource) toRequest(ctx context.Context, data *stageMutualTLSModel) (*api.MutualTLSStageRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	certificateAuthorities, d := helpers.SliceOrEmpty[string](ctx, data.CertificateAuthorities)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.MutualTLSStageRequest{
		Name:                   data.Name.ValueString(),
		Mode:                   api.StageModeEnum(data.Mode.ValueString()),
		CertAttribute:          api.CertAttributeEnum(data.CertAttribute.ValueString()),
		UserAttribute:          api.UserAttributeEnum(data.UserAttribute.ValueString()),
		CertificateAuthorities: certificateAuthorities,
	}, diags
}

func (r *stageMutualTLSResource) fromAPI(ctx context.Context, data *stageMutualTLSModel, res *api.MutualTLSStage) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(res.Pk)
	data.Name = types.StringValue(res.Name)
	data.Mode = types.StringValue(string(res.Mode))
	data.CertAttribute = types.StringValue(string(res.CertAttribute))
	data.UserAttribute = types.StringValue(string(res.UserAttribute))

	certificateAuthorities, d := helpers.MergeStringList(ctx, data.CertificateAuthorities, res.CertificateAuthorities)
	diags.Append(d...)
	data.CertificateAuthorities = certificateAuthorities

	return diags
}

func (r *stageMutualTLSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data stageMutualTLSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesMtlsCreate(ctx).MutualTLSStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageMutualTLSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data stageMutualTLSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesMtlsRetrieve(ctx, data.ID.ValueString()).Execute()
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

func (r *stageMutualTLSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data stageMutualTLSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.StagesAPI.StagesMtlsUpdate(ctx, data.ID.ValueString()).MutualTLSStageRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stageMutualTLSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data stageMutualTLSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.StagesAPI.StagesMtlsDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
