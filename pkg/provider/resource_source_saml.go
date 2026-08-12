package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &sourceSAMLResource{}
	_ resource.ResourceWithConfigure   = &sourceSAMLResource{}
	_ resource.ResourceWithImportState = &sourceSAMLResource{}
)

func newSourceSAMLResource() resource.Resource {
	return &sourceSAMLResource{}
}

type sourceSAMLResource struct {
	resourceBase
}

type sourceSAMLModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	UUID                     types.String `tfsdk:"uuid"`
	Slug                     types.String `tfsdk:"slug"`
	UserPathTemplate         types.String `tfsdk:"user_path_template"`
	AuthenticationFlow       types.String `tfsdk:"authentication_flow"`
	EnrollmentFlow           types.String `tfsdk:"enrollment_flow"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	Promoted                 types.Bool   `tfsdk:"promoted"`
	PolicyEngineMode         types.String `tfsdk:"policy_engine_mode"`
	UserMatchingMode         types.String `tfsdk:"user_matching_mode"`
	GroupMatchingMode        types.String `tfsdk:"group_matching_mode"`
	PreAuthenticationFlow    types.String `tfsdk:"pre_authentication_flow"`
	Issuer                   types.String `tfsdk:"issuer"`
	SSOURL                   types.String `tfsdk:"sso_url"`
	SLOURL                   types.String `tfsdk:"slo_url"`
	AllowIdpInitiated        types.Bool   `tfsdk:"allow_idp_initiated"`
	ForceAuthn               types.Bool   `tfsdk:"force_authn"`
	NameIDPolicy             types.String `tfsdk:"name_id_policy"`
	BindingType              types.String `tfsdk:"binding_type"`
	SigningKP                types.String `tfsdk:"signing_kp"`
	EncryptionKP             types.String `tfsdk:"encryption_kp"`
	VerificationKP           types.String `tfsdk:"verification_kp"`
	SignedAssertion          types.Bool   `tfsdk:"signed_assertion"`
	SignedResponse           types.Bool   `tfsdk:"signed_response"`
	DigestAlgorithm          types.String `tfsdk:"digest_algorithm"`
	SignatureAlgorithm       types.String `tfsdk:"signature_algorithm"`
	TemporaryUserDeleteAfter types.String `tfsdk:"temporary_user_delete_after"`
	Metadata                 types.String `tfsdk:"metadata"`
	PropertyMappings         types.List   `tfsdk:"property_mappings"`
	PropertyMappingsGroup    types.List   `tfsdk:"property_mappings_group"`
}

func (r *sourceSAMLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_saml"
}

func (r *sourceSAMLResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	userPathTemplateDefault := helpers.StringDefault("goauthentik.io/sources/%(slug)s")
	enabledDefault := helpers.BoolDefault(true)
	promotedDefault := helpers.BoolDefault(false)
	policyEngineModeDefault := helpers.StringDefault(string(api.POLICYENGINEMODE_ANY))
	userMatchingModeDefault := helpers.StringDefault(string(api.USERMATCHINGMODEENUM_IDENTIFIER))
	groupMatchingModeDefault := helpers.StringDefault(string(api.GROUPMATCHINGMODEENUM_IDENTIFIER))
	allowIdpInitiatedDefault := helpers.BoolDefault(false)
	forceAuthnDefault := helpers.BoolDefault(false)
	nameIDPolicyDefault := helpers.StringDefault(string(api.SAMLNAMEIDPOLICYENUM_URN_OASIS_NAMES_TC_SAML_2_0_NAMEID_FORMAT_PERSISTENT))
	bindingTypeDefault := helpers.StringDefault(string(api.BINDINGTYPEENUM_REDIRECT))
	signedAssertionDefault := helpers.BoolDefault(false)
	signedResponseDefault := helpers.BoolDefault(false)
	digestAlgorithmDefault := helpers.StringDefault(string(api.DIGESTALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLENCSHA256))
	signatureAlgorithmDefault := helpers.StringDefault(string(api.SIGNATUREALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLDSIG_MORERSA_SHA256))
	temporaryUserDeleteAfterDefault := helpers.StringDefault("days=1")

	resp.Schema = schema.Schema{
		MarkdownDescription: "Directory --- ",
		Attributes: map[string]schema.Attribute{
			// H3: id is res.Slug, so no UseStateForUnknown - see resource_source_scim.go.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"uuid": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: helpers.Desc("", helpers.Generated()),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Required: true,
			},
			"user_path_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userPathTemplateDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(userPathTemplateDefault.Value())),
			},
			"authentication_flow": schema.StringAttribute{
				Optional: true,
			},
			"enrollment_flow": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             enabledDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(enabledDefault.Value())),
			},
			"promoted": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             promotedDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(promotedDefault.Value())),
			},
			"policy_engine_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             policyEngineModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues), helpers.WithDefault(policyEngineModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedPolicyEngineModeEnumValues),
				},
			},
			"user_matching_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             userMatchingModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues), helpers.WithDefault(userMatchingModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedUserMatchingModeEnumEnumValues),
				},
			},
			"group_matching_mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             groupMatchingModeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues), helpers.WithDefault(groupMatchingModeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedGroupMatchingModeEnumEnumValues),
				},
			},
			"pre_authentication_flow": schema.StringAttribute{
				Required: true,
			},
			"issuer": schema.StringAttribute{
				Optional: true,
			},
			"sso_url": schema.StringAttribute{
				Required: true,
			},
			"slo_url": schema.StringAttribute{
				Optional: true,
			},
			"allow_idp_initiated": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             allowIdpInitiatedDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(allowIdpInitiatedDefault.Value())),
			},
			"force_authn": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             forceAuthnDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(forceAuthnDefault.Value())),
			},
			"name_id_policy": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             nameIDPolicyDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedSAMLNameIDPolicyEnumEnumValues), helpers.WithDefault(nameIDPolicyDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedSAMLNameIDPolicyEnumEnumValues),
				},
			},
			"binding_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             bindingTypeDefault,
				MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedBindingTypeEnumEnumValues), helpers.WithDefault(bindingTypeDefault.Value())),
				Validators: []validator.String{
					helpers.OneOf(api.AllowedBindingTypeEnumEnumValues),
				},
			},
			"signing_kp": schema.StringAttribute{
				Optional: true,
			},
			"encryption_kp": schema.StringAttribute{
				Optional: true,
			},
			"verification_kp": schema.StringAttribute{
				Optional: true,
			},
			"signed_assertion": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             signedAssertionDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(signedAssertionDefault.Value())),
			},
			"signed_response": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             signedResponseDefault,
				MarkdownDescription: helpers.Desc("", helpers.WithDefault(signedResponseDefault.Value())),
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
			"temporary_user_delete_after": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             temporaryUserDeleteAfterDefault,
				MarkdownDescription: helpers.Desc(helpers.RelativeDurationDescription, helpers.WithDefault(temporaryUserDeleteAfterDefault.Value())),
				Validators: []validator.String{
					helpers.RelativeDuration(),
				},
			},
			// Computed-only, and unlike every other attribute here it comes from a genuinely
			// separate endpoint (SourcesSamlMetadataRetrieve), not from the source object.
			// It embeds the slug, so no UseStateForUnknown.
			"metadata": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: helpers.Desc("SAML Metadata", helpers.Generated()),
			},
			"property_mappings": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"property_mappings_group": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *sourceSAMLResource) toRequest(ctx context.Context, data *sourceSAMLModel) (*api.SAMLSourceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	propertyMappings, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappings)
	diags.Append(d...)

	propertyMappingsGroup, d := helpers.SliceOrEmpty[string](ctx, data.PropertyMappingsGroup)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}

	return &api.SAMLSourceRequest{
		Name:                  data.Name.ValueString(),
		Slug:                  data.Slug.ValueString(),
		Enabled:               new(data.Enabled.ValueBool()),
		Promoted:              new(data.Promoted.ValueBool()),
		UserPathTemplate:      new(data.UserPathTemplate.ValueString()),
		PolicyEngineMode:      api.PolicyEngineMode(data.PolicyEngineMode.ValueString()).Ptr(),
		UserMatchingMode:      api.UserMatchingModeEnum(data.UserMatchingMode.ValueString()).Ptr(),
		GroupMatchingMode:     api.GroupMatchingModeEnum(data.GroupMatchingMode.ValueString()).Ptr(),
		AuthenticationFlow:    *api.NewNullableString(helpers.StringPtr(data.AuthenticationFlow)),
		EnrollmentFlow:        *api.NewNullableString(helpers.StringPtr(data.EnrollmentFlow)),
		PreAuthenticationFlow: data.PreAuthenticationFlow.ValueString(),
		UserPropertyMappings:  propertyMappings,
		GroupPropertyMappings: propertyMappingsGroup,
		SsoUrl:                data.SSOURL.ValueString(),
		SloUrl:                *api.NewNullableString(helpers.StringPtr(data.SLOURL)),
		SigningKp:             *api.NewNullableString(helpers.StringPtr(data.SigningKP)),
		EncryptionKp:          *api.NewNullableString(helpers.StringPtr(data.EncryptionKP)),
		VerificationKp:        *api.NewNullableString(helpers.StringPtr(data.VerificationKP)),
		SignedAssertion:       new(data.SignedAssertion.ValueBool()),
		SignedResponse:        new(data.SignedResponse.ValueBool()),
		// SDKv2 sent issuer via new(d.Get(...).(string)), i.e. always a pointer and "" when
		// unset, so it keeps StringPtrEmpty semantics rather than being omitted.
		Issuer:                   new(data.Issuer.ValueString()),
		AllowIdpInitiated:        new(data.AllowIdpInitiated.ValueBool()),
		ForceAuthn:               new(data.ForceAuthn.ValueBool()),
		TemporaryUserDeleteAfter: new(data.TemporaryUserDeleteAfter.ValueString()),
		BindingType:              api.BindingTypeEnum(data.BindingType.ValueString()).Ptr(),
		DigestAlgorithm:          api.DigestAlgorithmEnum(data.DigestAlgorithm.ValueString()).Ptr(),
		SignatureAlgorithm:       api.SignatureAlgorithmEnum(data.SignatureAlgorithm.ValueString()).Ptr(),
		NameIdPolicy:             api.SAMLNameIDPolicyEnum(data.NameIDPolicy.ValueString()).Ptr(),
	}, diags
}

// fromAPI does not set Metadata: it lives on a different endpoint, so the callers fetch it
// separately and assign it after this returns.
func (r *sourceSAMLResource) fromAPI(ctx context.Context, data *sourceSAMLModel, res *api.SAMLSource) diag.Diagnostics {
	var diags diag.Diagnostics

	// H3: id is the slug, not the Pk.
	data.ID = types.StringValue(res.Slug)
	data.Name = types.StringValue(res.Name)
	data.Slug = types.StringValue(res.Slug)
	data.UUID = types.StringValue(res.Pk)

	propertyMappings, d := helpers.MergeStringList(ctx, data.PropertyMappings, res.UserPropertyMappings)
	diags.Append(d...)
	data.PropertyMappings = propertyMappings

	propertyMappingsGroup, d := helpers.MergeStringList(ctx, data.PropertyMappingsGroup, res.GroupPropertyMappings)
	diags.Append(d...)
	data.PropertyMappingsGroup = propertyMappingsGroup

	// Nullable API fields.
	data.AuthenticationFlow = helpers.StringPtrOrNull(res.AuthenticationFlow.Get())
	data.EnrollmentFlow = helpers.StringPtrOrNull(res.EnrollmentFlow.Get())
	data.SLOURL = helpers.StringPtrOrNull(res.SloUrl.Get())
	data.SigningKP = helpers.StringPtrOrNull(res.SigningKp.Get())
	data.EncryptionKP = helpers.StringPtrOrNull(res.EncryptionKp.Get())
	data.VerificationKP = helpers.StringPtrOrNull(res.VerificationKp.Get())
	// No Default, so prior-aware.
	data.Issuer = helpers.StringOrNull(data.Issuer, res.GetIssuer())
	// Required, so never null.
	data.PreAuthenticationFlow = types.StringValue(res.PreAuthenticationFlow)
	data.SSOURL = types.StringValue(res.SsoUrl)
	// Everything below has a Default and takes the API value verbatim.
	data.UserPathTemplate = types.StringValue(res.GetUserPathTemplate())
	data.Enabled = types.BoolValue(res.GetEnabled())
	data.Promoted = types.BoolValue(res.GetPromoted())
	data.PolicyEngineMode = types.StringValue(string(res.GetPolicyEngineMode()))
	data.UserMatchingMode = types.StringValue(string(res.GetUserMatchingMode()))
	data.GroupMatchingMode = types.StringValue(string(res.GetGroupMatchingMode()))
	data.AllowIdpInitiated = types.BoolValue(res.GetAllowIdpInitiated())
	data.ForceAuthn = types.BoolValue(res.GetForceAuthn())
	data.NameIDPolicy = types.StringValue(string(res.GetNameIdPolicy()))
	data.BindingType = types.StringValue(string(res.GetBindingType()))
	data.SignedAssertion = types.BoolValue(res.GetSignedAssertion())
	data.SignedResponse = types.BoolValue(res.GetSignedResponse())
	data.DigestAlgorithm = types.StringValue(string(res.GetDigestAlgorithm()))
	data.SignatureAlgorithm = types.StringValue(string(res.GetSignatureAlgorithm()))
	data.TemporaryUserDeleteAfter = types.StringValue(res.GetTemporaryUserDeleteAfter())

	return diags
}

// readMetadata fills the Computed metadata attribute from its own endpoint.
func (r *sourceSAMLResource) readMetadata(ctx context.Context, data *sourceSAMLModel) diag.Diagnostics {
	var diags diag.Diagnostics

	meta, hr, err := r.client.SourcesAPI.SourcesSamlMetadataRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		diags.Append(helpers.HTTPError(hr, err)...)
		return diags
	}
	data.Metadata = types.StringValue(meta.GetMetadata())
	return diags
}

func (r *sourceSAMLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer r.span(ctx, "create")()

	var data sourceSAMLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesSamlCreate(ctx).SAMLSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.readMetadata(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceSAMLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer r.span(ctx, "read")()

	var data sourceSAMLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesSamlRetrieve(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		if helpers.IsNotFound(hr) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.readMetadata(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceSAMLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer r.span(ctx, "update")()

	var data sourceSAMLModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// H3/discovery #4: the plan's id is unknown when the slug changes.
	var priorID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &priorID)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body, diags := r.toRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, hr, err := r.client.SourcesAPI.SourcesSamlUpdate(ctx, priorID.ValueString()).SAMLSourceRequest(*body).Execute()
	if err != nil {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
		return
	}

	resp.Diagnostics.Append(r.fromAPI(ctx, &data, res)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(r.readMetadata(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sourceSAMLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer r.span(ctx, "delete")()

	var data sourceSAMLModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hr, err := r.client.SourcesAPI.SourcesSamlDestroy(ctx, data.ID.ValueString()).Execute()
	if err != nil && !helpers.IsNotFound(hr) {
		resp.Diagnostics.Append(helpers.HTTPError(hr, err)...)
	}
}
