package provider

import (
	"context"
	"strconv"

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
	_ resource.Resource                = &providerWSFederationResource{}
	_ resource.ResourceWithConfigure   = &providerWSFederationResource{}
	_ resource.ResourceWithImportState = &providerWSFederationResource{}
)

func newProviderWSFederationResource() resource.Resource {
	return &providerWSFederationResource{}
}

type providerWSFederationResource struct {
	resourceBase
}

type providerWSFederationModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	AuthenticationFlow          types.String `tfsdk:"authentication_flow"`
	AuthorizationFlow           types.String `tfsdk:"authorization_flow"`
	InvalidationFlow            types.String `tfsdk:"invalidation_flow"`
	ReplyURL                    types.String `tfsdk:"reply_url"`
	Wtrealm                     types.String `tfsdk:"wtrealm"`
	PropertyMappings            types.List   `tfsdk:"property_mappings"`
	AssertionValidNotBefore     types.String `tfsdk:"assertion_valid_not_before"`
	AssertionValidNotOnOrAfter  types.String `tfsdk:"assertion_valid_not_on_or_after"`
	SessionValidNotOnOrAfter    types.String `tfsdk:"session_valid_not_on_or_after"`
	NameIDMapping               types.String `tfsdk:"name_id_mapping"`
	AuthnContextClassRefMapping types.String `tfsdk:"authn_context_class_ref_mapping"`
	DigestAlgorithm             types.String `tfsdk:"digest_algorithm"`
	SignatureAlgorithm          types.String `tfsdk:"signature_algorithm"`
	SigningKP                   types.String `tfsdk:"signing_kp"`
	SignAssertion               types.Bool   `tfsdk:"sign_assertion"`
	EncryptionKP                types.String `tfsdk:"encryption_kp"`
	SignLogoutRequest           types.Bool   `tfsdk:"sign_logout_request"`
}

func (r *providerWSFederationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_ws_federation"
}

func (r *providerWSFederationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	assertionValidNotBeforeDefault := helpers.StringDefault("minutes=-5")
	assertionValidNotOnOrAfterDefault := helpers.StringDefault("minutes=5")
	sessionValidNotOnOrAfterDefault := helpers.StringDefault("minutes=86400")
	digestAlgorithmDefault := helpers.StringDefault(string(api.DIGESTALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLENCSHA256))
	signatureAlgorithmDefault := helpers.StringDefault(string(api.SIGNATUREALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLDSIG_MORERSA_SHA256))
	signAssertionDefault := helpers.BoolDefault(true)
	signLogoutRequestDefault := helpers.BoolDefault(false)

	resp.Schema = schema.Schema{
		MarkdownDescription: "Applications --- ",
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
			"authentication_flow": schema.StringAttribute{
				Optional: true,
			},
			"authorization_flow": schema.StringAttribute{
				Required: true,
			},
			"invalidation_flow": schema.StringAttribute{
				Required: true,
			},
			"reply_url": schema.StringAttribute{
				Required: true,
			},
			"wtrealm": schema.StringAttribute{
				Required: true,
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"assertion_valid_not_before": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             assertionValidNotBeforeDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(assertionValidNotBeforeDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"assertion_valid_not_on_or_after": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             assertionValidNotOnOrAfterDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(assertionValidNotOnOrAfterDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"session_valid_not_on_or_after": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             sessionValidNotOnOrAfterDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(sessionValidNotOnOrAfterDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			"name_id_mapping": schema.StringAttribute{
				Optional: true,
			},
			"authn_context_class_ref_mapping": schema.StringAttribute{
				Optional: true,
			},
			"digest_algorithm": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             digestAlgorithmDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedDigestAlgorithmEnumEnumValues), helpers.WithDefault(digestAlgorithmDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedDigestAlgorithmEnumEnumValues),
				},
			},
			"signature_algorithm": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             signatureAlgorithmDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedSignatureAlgorithmEnumEnumValues), helpers.WithDefault(signatureAlgorithmDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedSignatureAlgorithmEnumEnumValues),
				},
			},
			"signing_kp": schema.StringAttribute{
				Optional: true,
			},
			"sign_assertion": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             signAssertionDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(signAssertionDefault.Value())),
			},
			"encryption_kp": schema.StringAttribute{
				Optional: true,
			},
			"sign_logout_request": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             signLogoutRequestDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(signLogoutRequestDefault.Value())),
			},
		},
	}
}

func (r *providerWSFederationResource) toRequest(ctx context.Context, data *providerWSFederationModel) (*api.WSFederationProviderRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	return &api.WSFederationProviderRequest{
		Name:                        data.Name.ValueString(),
		AuthorizationFlow:           data.AuthorizationFlow.ValueString(),
		InvalidationFlow:            data.InvalidationFlow.ValueString(),
		ReplyUrl:                    data.ReplyURL.ValueString(),
		Wtrealm:                     data.Wtrealm.ValueString(),
		AssertionValidNotBefore:     new(data.AssertionValidNotBefore.ValueString()),
		AssertionValidNotOnOrAfter:  new(data.AssertionValidNotOnOrAfter.ValueString()),
		SessionValidNotOnOrAfter:    new(data.SessionValidNotOnOrAfter.ValueString()),
		DigestAlgorithm:             api.DigestAlgorithmEnum(data.DigestAlgorithm.ValueString()).Ptr(),
		SignatureAlgorithm:          api.SignatureAlgorithmEnum(data.SignatureAlgorithm.ValueString()).Ptr(),
		PropertyMappings:            propertyMappings,
		SignAssertion:               new(data.SignAssertion.ValueBool()),
		AuthenticationFlow:          *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		NameIdMapping:               *api.NewNullableString(helpers.StringPtr(data.NameIDMapping)),
		AuthnContextClassRefMapping: *api.NewNullableString(helpers.StringPtr(data.AuthnContextClassRefMapping)),
		EncryptionKp:                *api.NewNullableString(helpers.StringPtr(data.EncryptionKP)),
		SigningKp:                   *api.NewNullableString(helpers.StringPtr(data.SigningKP)),
		// Another GetP[bool] casualty: SDKv2 omitted sign_logout_request whenever it was
		// false, so - the API leaving absent fields unchanged on PUT - it could not be turned
		// back off once enabled. It has a Default, so it is never null here and the concrete
		// value is always sent. Sixth instance of this defect across the migration.
		SignLogoutRequest: new(data.SignLogoutRequest.ValueBool()),
	}, diags
}

func (r *providerWSFederationResource) fromAPI(ctx context.Context, data *providerWSFederationModel, res *api.WSFederationProvider) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(strconv.Itoa(int(res.Pk)))
	data.Name = types.StringValue(res.Name)
	data.AuthorizationFlow = types.StringValue(res.AuthorizationFlow)
	data.InvalidationFlow = types.StringValue(res.InvalidationFlow)
	data.ReplyURL = types.StringValue(res.ReplyUrl)
	data.Wtrealm = types.StringValue(res.Wtrealm)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.PropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.NameIDMapping = helpers.StringPtrOrNull(res.NameIdMapping.Get())
	data.AuthnContextClassRefMapping = helpers.StringPtrOrNull(res.AuthnContextClassRefMapping.Get())
	data.SigningKP = helpers.StringPtrOrNull(res.SigningKp.Get())
	data.EncryptionKP = helpers.StringPtrOrNull(res.EncryptionKp.Get())
	// Have Defaults, so verbatim.
	data.AssertionValidNotBefore = types.StringValue(res.GetAssertionValidNotBefore())
	data.AssertionValidNotOnOrAfter = types.StringValue(res.GetAssertionValidNotOnOrAfter())
	data.SessionValidNotOnOrAfter = types.StringValue(res.GetSessionValidNotOnOrAfter())
	data.DigestAlgorithm = types.StringValue(string(res.GetDigestAlgorithm()))
	data.SignatureAlgorithm = types.StringValue(string(res.GetSignatureAlgorithm()))
	data.SignAssertion = types.BoolValue(res.GetSignAssertion())
	data.SignLogoutRequest = types.BoolValue(res.GetSignLogoutRequest())

	return diags
}

func (r *providerWSFederationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data providerWSFederationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersWsfedCreate(ctx).WSFederationProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerWSFederationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data providerWSFederationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersWsfedRetrieve(ctx, id).Execute()
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

func (r *providerWSFederationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data providerWSFederationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.ProvidersAPI.ProvidersWsfedUpdate(ctx, id).WSFederationProviderRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *providerWSFederationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data providerWSFederationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, diags := helpers.ParseInt32ID(data.ID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.ProvidersAPI.ProvidersWsfedDestroy(ctx, id).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
