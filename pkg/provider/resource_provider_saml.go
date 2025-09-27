package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
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
				Default:          api.SPBINDINGENUM_REDIRECT,
				Description:      helpers.EnumToDescription(api.AllowedSpBindingEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSpBindingEnumEnumValues),
			},
			"default_relay_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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
		Audience:                    api.PtrString(d.Get("audience").(string)),
		Issuer:                      api.PtrString(d.Get("issuer").(string)),
		AssertionValidNotBefore:     api.PtrString(d.Get("assertion_valid_not_before").(string)),
		AssertionValidNotOnOrAfter:  api.PtrString(d.Get("assertion_valid_not_on_or_after").(string)),
		SessionValidNotOnOrAfter:    api.PtrString(d.Get("session_valid_not_on_or_after").(string)),
		DigestAlgorithm:             api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string)).Ptr(),
		SignatureAlgorithm:          api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string)).Ptr(),
		SpBinding:                   api.SpBindingEnum(d.Get("sp_binding").(string)).Ptr(),
		PropertyMappings:            castSlice[string](d.Get("property_mappings").([]interface{})),
		SignAssertion:               api.PtrBool(d.Get("sign_assertion").(bool)),
		SignResponse:                api.PtrBool(d.Get("sign_response").(bool)),
		AuthenticationFlow:          *api.NewNullableString(getP[string](d, "authentication_flow")),
		NameIdMapping:               *api.NewNullableString(getP[string](d, "name_id_mapping")),
		AuthnContextClassRefMapping: *api.NewNullableString(getP[string](d, "authn_context_class_ref_mapping")),
		EncryptionKp:                *api.NewNullableString(getP[string](d, "encryption_kp")),
		SigningKp:                   *api.NewNullableString(getP[string](d, "signing_kp")),
		VerificationKp:              *api.NewNullableString(getP[string](d, "verification_kp")),
		DefaultRelayState:           getP[string](d, "default_relay_state"),
	}
	return &r
}

func resourceProviderSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersSamlCreate(ctx).SAMLProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSAMLRead(ctx, d, m)
}

func resourceProviderSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setWrapper(d, "name", res.Name)
	setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	setWrapper(d, "authorization_flow", res.AuthorizationFlow)
	setWrapper(d, "invalidation_flow", res.InvalidationFlow)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", helpers.ListConsistentMerge(localMappings, res.PropertyMappings))

	setWrapper(d, "acs_url", res.AcsUrl)
	setWrapper(d, "audience", res.Audience)
	setWrapper(d, "issuer", res.Issuer)
	setWrapper(d, "sp_binding", res.SpBinding)
	setWrapper(d, "assertion_valid_not_before", res.AssertionValidNotBefore)
	setWrapper(d, "assertion_valid_not_on_or_after", res.AssertionValidNotOnOrAfter)
	setWrapper(d, "session_valid_not_on_or_after", res.SessionValidNotOnOrAfter)
	setWrapper(d, "sign_assertion", res.SignAssertion)
	setWrapper(d, "sign_response", res.SignResponse)
	if res.NameIdMapping.IsSet() {
		setWrapper(d, "name_id_mapping", res.NameIdMapping.Get())
	}
	if res.AuthnContextClassRefMapping.IsSet() {
		setWrapper(d, "authn_context_class_ref_mapping", res.AuthnContextClassRefMapping.Get())
	}
	if res.SigningKp.IsSet() {
		setWrapper(d, "signing_kp", res.SigningKp.Get())
	}
	if res.VerificationKp.IsSet() {
		setWrapper(d, "verification_kp", res.VerificationKp.Get())
	}
	if res.EncryptionKp.IsSet() {
		setWrapper(d, "encryption_kp", res.EncryptionKp.Get())
	}
	setWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	setWrapper(d, "signature_algorithm", res.SignatureAlgorithm)
	setWrapper(d, "default_relay_state", res.DefaultRelayState)

	setWrapper(d, "url_sso_init", res.UrlSsoInit)
	setWrapper(d, "url_sso_post", res.UrlSsoPost)
	setWrapper(d, "url_sso_redirect", res.UrlSsoRedirect)
	setWrapper(d, "url_slo_post", res.UrlSloPost)
	setWrapper(d, "url_slo_redirect", res.UrlSloRedirect)
	return diags
}

func resourceProviderSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceProviderSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
