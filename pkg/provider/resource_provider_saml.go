package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderSAMLCreate,
		ReadContext:   resourceProviderSAMLRead,
		UpdateContext: resourceProviderSAMLUpdate,
		DeleteContext: resourceProviderSAMLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"url_sso_init": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"url_sso_post": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"url_sso_redirect": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"url_slo_post": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"url_slo_redirect": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"invalidation_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"acs_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"audience": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "authentik",
			},
			"assertion_valid_not_before": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=-5",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"assertion_valid_not_on_or_after": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=5",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"session_valid_not_on_or_after": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=86400",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"name_id_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authn_context_class_ref_mapping": {
				Type:     schema.TypeString,
				Optional: true,
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
			"signing_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sign_assertion": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"sign_response": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"verification_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sp_binding": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SAMLBINDINGSENUM_REDIRECT,
				Description:      helpers.EnumToDescription(api.AllowedSAMLBindingsEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSAMLBindingsEnumEnumValues),
			},
			"default_relay_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"sls_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sign_logout_request": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sls_binding": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SAMLBINDINGSENUM_REDIRECT,
				Description:      helpers.EnumToDescription(api.AllowedSAMLBindingsEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSAMLBindingsEnumEnumValues),
			},
			"logout_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SAMLPROVIDERLOGOUTMETHODENUM_FRONTCHANNEL_IFRAME,
				Description:      helpers.EnumToDescription(api.AllowedSAMLProviderLogoutMethodEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSAMLProviderLogoutMethodEnumEnumValues),
			},
		},
	}
}

func resourceProviderSAMLSchemaToProvider(d *schema.ResourceData) *api.SAMLProviderRequest {
	r := api.SAMLProviderRequest{
		Name:                        d.Get("name").(string),
		AuthorizationFlow:           d.Get("authorization_flow").(string),
		InvalidationFlow:            d.Get("invalidation_flow").(string),
		AcsUrl:                      d.Get("acs_url").(string),
		Audience:                    new(d.Get("audience").(string)),
		Issuer:                      new(d.Get("issuer").(string)),
		AssertionValidNotBefore:     new(d.Get("assertion_valid_not_before").(string)),
		AssertionValidNotOnOrAfter:  new(d.Get("assertion_valid_not_on_or_after").(string)),
		SessionValidNotOnOrAfter:    new(d.Get("session_valid_not_on_or_after").(string)),
		DigestAlgorithm:             api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string)).Ptr(),
		SignatureAlgorithm:          api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string)).Ptr(),
		SpBinding:                   api.SAMLBindingsEnum(d.Get("sp_binding").(string)).Ptr(),
		PropertyMappings:            helpers.CastSlice[string](d, "property_mappings"),
		SignAssertion:               new(d.Get("sign_assertion").(bool)),
		SignResponse:                new(d.Get("sign_response").(bool)),
		AuthenticationFlow:          *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		NameIdMapping:               *api.NewNullableString(helpers.GetP[string](d, "name_id_mapping")),
		AuthnContextClassRefMapping: *api.NewNullableString(helpers.GetP[string](d, "authn_context_class_ref_mapping")),
		EncryptionKp:                *api.NewNullableString(helpers.GetP[string](d, "encryption_kp")),
		SigningKp:                   *api.NewNullableString(helpers.GetP[string](d, "signing_kp")),
		VerificationKp:              *api.NewNullableString(helpers.GetP[string](d, "verification_kp")),
		DefaultRelayState:           helpers.GetP[string](d, "default_relay_state"),
		SlsUrl:                      helpers.GetP[string](d, "sls_url"),
		SignLogoutRequest:           helpers.GetP[bool](d, "sign_logout_request"),
		SlsBinding:                  helpers.CastString[api.SAMLBindingsEnum](helpers.GetP[string](d, "sls_binding")),
		LogoutMethod:                helpers.CastString[api.SAMLProviderLogoutMethodEnum](helpers.GetP[string](d, "logout_method")),
	}
	return &r
}

func resourceProviderSAMLCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersSamlCreate(ctx).SAMLProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSAMLRead(ctx, d, m)
}

func resourceProviderSAMLRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersSamlRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "invalidation_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))

	helpers.SetWrapper(d, "acs_url", res.AcsUrl)
	helpers.SetWrapper(d, "audience", res.Audience)
	helpers.SetWrapper(d, "issuer", res.Issuer)
	helpers.SetWrapper(d, "sp_binding", res.SpBinding)
	helpers.SetWrapper(d, "assertion_valid_not_before", res.AssertionValidNotBefore)
	helpers.SetWrapper(d, "assertion_valid_not_on_or_after", res.AssertionValidNotOnOrAfter)
	helpers.SetWrapper(d, "session_valid_not_on_or_after", res.SessionValidNotOnOrAfter)
	helpers.SetWrapper(d, "sign_assertion", res.SignAssertion)
	helpers.SetWrapper(d, "sign_response", res.SignResponse)
	helpers.SetWrapper(d, "name_id_mapping", res.NameIdMapping.Get())
	helpers.SetWrapper(d, "authn_context_class_ref_mapping", res.AuthnContextClassRefMapping.Get())
	helpers.SetWrapper(d, "signing_kp", res.SigningKp.Get())
	helpers.SetWrapper(d, "verification_kp", res.VerificationKp.Get())
	helpers.SetWrapper(d, "encryption_kp", res.EncryptionKp.Get())
	helpers.SetWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	helpers.SetWrapper(d, "signature_algorithm", res.SignatureAlgorithm)
	helpers.SetWrapper(d, "default_relay_state", res.DefaultRelayState)

	helpers.SetWrapper(d, "url_sso_init", res.UrlSsoInit)
	helpers.SetWrapper(d, "url_sso_post", res.UrlSsoPost)
	helpers.SetWrapper(d, "url_sso_redirect", res.UrlSsoRedirect)
	helpers.SetWrapper(d, "url_slo_post", res.UrlSloPost)
	helpers.SetWrapper(d, "url_slo_redirect", res.UrlSloRedirect)
	helpers.SetWrapper(d, "sls_url", res.SlsUrl)
	helpers.SetWrapper(d, "sign_logout_request", res.SignLogoutRequest)
	helpers.SetWrapper(d, "sls_binding", res.SlsBinding)
	helpers.SetWrapper(d, "logout_method", res.LogoutMethod)
	return diags
}

func resourceProviderSAMLUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersSamlUpdate(ctx, int32(id)).SAMLProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSAMLRead(ctx, d, m)
}

func resourceProviderSAMLDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersSamlDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
