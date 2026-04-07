package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderWSFederation() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderWSFederationCreate,
		ReadContext:   resourceProviderWSFederationRead,
		UpdateContext: resourceProviderWSFederationUpdate,
		DeleteContext: resourceProviderWSFederationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"reply_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wtrealm": {
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
				Default:          api.DIGESTALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLENCSHA256,
				Description:      helpers.EnumToDescription(api.AllowedDigestAlgorithmEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedDigestAlgorithmEnumEnumValues),
			},
			"signature_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SIGNATUREALGORITHMENUM_HTTP___WWW_W3_ORG_2001_04_XMLDSIG_MORERSA_SHA256,
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
			"encryption_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sign_logout_request": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceProviderWSFederationSchemaToProvider(d *schema.ResourceData) *api.WSFederationProviderRequest {
	r := api.WSFederationProviderRequest{
		Name:                        d.Get("name").(string),
		AuthorizationFlow:           d.Get("authorization_flow").(string),
		InvalidationFlow:            d.Get("invalidation_flow").(string),
		ReplyUrl:                    d.Get("reply_url").(string),
		Wtrealm:                     d.Get("wtrealm").(string),
		AssertionValidNotBefore:     new(d.Get("assertion_valid_not_before").(string)),
		AssertionValidNotOnOrAfter:  new(d.Get("assertion_valid_not_on_or_after").(string)),
		SessionValidNotOnOrAfter:    new(d.Get("session_valid_not_on_or_after").(string)),
		DigestAlgorithm:             api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string)).Ptr(),
		SignatureAlgorithm:          api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string)).Ptr(),
		PropertyMappings:            helpers.CastSlice[string](d, "property_mappings"),
		SignAssertion:               new(d.Get("sign_assertion").(bool)),
		AuthenticationFlow:          *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		NameIdMapping:               *api.NewNullableString(helpers.GetP[string](d, "name_id_mapping")),
		AuthnContextClassRefMapping: *api.NewNullableString(helpers.GetP[string](d, "authn_context_class_ref_mapping")),
		EncryptionKp:                *api.NewNullableString(helpers.GetP[string](d, "encryption_kp")),
		SigningKp:                   *api.NewNullableString(helpers.GetP[string](d, "signing_kp")),
		SignLogoutRequest:           helpers.GetP[bool](d, "sign_logout_request"),
	}
	return &r
}

func resourceProviderWSFederationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderWSFederationSchemaToProvider(d)

	res, hr, err := c.client.ProvidersAPI.ProvidersWsfedCreate(ctx).WSFederationProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderWSFederationRead(ctx, d, m)
}

func resourceProviderWSFederationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersAPI.ProvidersWsfedRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "invalidation_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "reply_url", res.ReplyUrl)
	helpers.SetWrapper(d, "wtrealm", res.Wtrealm)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))

	helpers.SetWrapper(d, "assertion_valid_not_before", res.AssertionValidNotBefore)
	helpers.SetWrapper(d, "assertion_valid_not_on_or_after", res.AssertionValidNotOnOrAfter)
	helpers.SetWrapper(d, "session_valid_not_on_or_after", res.SessionValidNotOnOrAfter)
	helpers.SetWrapper(d, "sign_assertion", res.SignAssertion)
	helpers.SetWrapper(d, "name_id_mapping", res.NameIdMapping.Get())
	helpers.SetWrapper(d, "authn_context_class_ref_mapping", res.AuthnContextClassRefMapping.Get())
	helpers.SetWrapper(d, "signing_kp", res.SigningKp.Get())
	helpers.SetWrapper(d, "encryption_kp", res.EncryptionKp.Get())
	helpers.SetWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	helpers.SetWrapper(d, "signature_algorithm", res.SignatureAlgorithm)

	helpers.SetWrapper(d, "sign_logout_request", res.SignLogoutRequest)
	return diags
}

func resourceProviderWSFederationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderWSFederationSchemaToProvider(d)

	res, hr, err := c.client.ProvidersAPI.ProvidersWsfedUpdate(ctx, int32(id)).WSFederationProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderWSFederationRead(ctx, d, m)
}

func resourceProviderWSFederationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersAPI.ProvidersWsfedDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
