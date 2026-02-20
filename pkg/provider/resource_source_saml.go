package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceSourceSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceSAMLCreate,
		ReadContext:   resourceSourceSAMLRead,
		UpdateContext: resourceSourceSAMLUpdate,
		DeleteContext: resourceSourceSAMLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "goauthentik.io/sources/%(slug)s",
				Optional: true,
			},
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"promoted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"user_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserMatchingModeEnumEnumValues),
			},
			"group_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.GROUPMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedGroupMatchingModeEnumEnumValues),
			},

			"pre_authentication_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sso_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slo_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_idp_initiated": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"name_id_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SAMLNAMEIDPOLICYENUM__2_0NAMEID_FORMATPERSISTENT,
				Description:      helpers.EnumToDescription(api.AllowedSAMLNameIDPolicyEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSAMLNameIDPolicyEnumEnumValues),
			},
			"binding_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.BINDINGTYPEENUM_REDIRECT,
				Description:      helpers.EnumToDescription(api.AllowedBindingTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedBindingTypeEnumEnumValues),
			},
			"signing_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verification_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signed_assertion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"signed_response": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"digest_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.DIGESTALGORITHMENUM__2001_04_XMLENCSHA256,
				Description:      helpers.EnumToDescription(api.AllowedDigestAlgorithmEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedDigestAlgorithmEnumEnumValues),
			},
			"signature_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SIGNATUREALGORITHMENUM__2001_04_XMLDSIG_MORERSA_SHA256,
				Description:      helpers.EnumToDescription(api.AllowedSignatureAlgorithmEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSignatureAlgorithmEnumEnumValues),
			},
			"temporary_user_delete_after": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "days=1",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},

			"metadata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML Metadata",
			},

			"property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property_mappings_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSourceSAMLSchemaToSource(d *schema.ResourceData) *api.SAMLSourceRequest {
	r := api.SAMLSourceRequest{
		Name:              d.Get("name").(string),
		Slug:              d.Get("slug").(string),
		Enabled:           new(d.Get("enabled").(bool)),
		Promoted:          new(d.Get("promoted").(bool)),
		UserPathTemplate:  new(d.Get("user_path_template").(string)),
		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:  api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode: api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),

		AuthenticationFlow:    *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		EnrollmentFlow:        *api.NewNullableString(helpers.GetP[string](d, "enrollment_flow")),
		PreAuthenticationFlow: d.Get("pre_authentication_flow").(string),
		UserPropertyMappings:  helpers.CastSlice[string](d, "property_mappings"),
		GroupPropertyMappings: helpers.CastSlice[string](d, "property_mappings_group"),

		SsoUrl:                   d.Get("sso_url").(string),
		SloUrl:                   *api.NewNullableString(helpers.GetP[string](d, "slo_url")),
		SigningKp:                *api.NewNullableString(helpers.GetP[string](d, "signing_kp")),
		EncryptionKp:             *api.NewNullableString(helpers.GetP[string](d, "encryption_kp")),
		VerificationKp:           *api.NewNullableString(helpers.GetP[string](d, "verification_kp")),
		SignedAssertion:          new(d.Get("signed_assertion").(bool)),
		SignedResponse:           new(d.Get("signed_response").(bool)),
		Issuer:                   new(d.Get("issuer").(string)),
		AllowIdpInitiated:        new(d.Get("allow_idp_initiated").(bool)),
		TemporaryUserDeleteAfter: new(d.Get("temporary_user_delete_after").(string)),
		BindingType:              api.BindingTypeEnum(d.Get("binding_type").(string)).Ptr(),
		DigestAlgorithm:          api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string)).Ptr(),
		SignatureAlgorithm:       api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string)).Ptr(),
		NameIdPolicy:             api.SAMLNameIDPolicyEnum(d.Get("name_id_policy").(string)).Ptr(),
	}
	return &r
}

func resourceSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlCreate(ctx).SAMLSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "user_path_template", res.UserPathTemplate)

	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "promoted", res.Promoted)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "user_matching_mode", res.UserMatchingMode)
	helpers.SetWrapper(d, "group_matching_mode", res.GroupMatchingMode)

	helpers.SetWrapper(d, "pre_authentication_flow", res.PreAuthenticationFlow)
	helpers.SetWrapper(d, "issuer", res.Issuer)
	helpers.SetWrapper(d, "sso_url", res.SsoUrl)
	helpers.SetWrapper(d, "slo_url", res.SloUrl.Get())
	helpers.SetWrapper(d, "allow_idp_initiated", res.AllowIdpInitiated)
	helpers.SetWrapper(d, "name_id_policy", res.NameIdPolicy)
	helpers.SetWrapper(d, "binding_type", res.BindingType)
	helpers.SetWrapper(d, "signing_kp", res.SigningKp.Get())
	helpers.SetWrapper(d, "encryption_kp", res.EncryptionKp.Get())
	helpers.SetWrapper(d, "verification_kp", res.VerificationKp.Get())
	helpers.SetWrapper(d, "signed_assertion", res.SignedAssertion)
	helpers.SetWrapper(d, "signed_response", res.SignedResponse)
	helpers.SetWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	helpers.SetWrapper(d, "signature_algorithm", res.SignatureAlgorithm)
	helpers.SetWrapper(d, "temporary_user_delete_after", res.TemporaryUserDeleteAfter)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.UserPropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.GroupPropertyMappings,
	))

	meta, hr, err := c.client.SourcesApi.SourcesSamlMetadataRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "metadata", meta.Metadata)
	return diags
}

func resourceSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlUpdate(ctx, d.Id()).SAMLSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
