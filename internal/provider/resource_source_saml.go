package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceSAML() *schema.Resource {
	return &schema.Resource{
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
			"authentication_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},
			"user_matching_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.USERMATCHINGMODEENUM_IDENTIFIER,
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.NAMEIDPOLICYENUM__2_0NAMEID_FORMATPERSISTENT,
			},
			"binding_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.BINDINGTYPEENUM_REDIRECT,
			},
			"signing_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digest_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.DIGESTALGORITHMENUM__2001_04_XMLENCSHA256,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.SIGNATUREALGORITHMENUM__2001_04_XMLDSIG_MORERSA_SHA256,
			},
			"temporary_user_delete_after": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=1",
			},
		},
	}
}

func resourceSourceSAMLSchemaToSource(d *schema.ResourceData) *api.SAMLSourceRequest {
	r := api.SAMLSourceRequest{
		Name:    d.Get("name").(string),
		Slug:    d.Get("slug").(string),
		Enabled: boolToPointer(d.Get("enabled").(bool)),

		PreAuthenticationFlow: d.Get("pre_authentication_flow").(string),

		SsoUrl:                   d.Get("sso_url").(string),
		Issuer:                   stringToPointer(d.Get("issuer").(string)),
		AllowIdpInitiated:        boolToPointer(d.Get("allow_idp_initiated").(bool)),
		TemporaryUserDeleteAfter: stringToPointer(d.Get("temporary_user_delete_after").(string)),
	}

	r.AuthenticationFlow.Set(stringToPointer(d.Get("authentication_flow").(string)))
	r.EnrollmentFlow.Set(stringToPointer(d.Get("enrollment_flow").(string)))

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	r.PolicyEngineMode = &pm

	bt := api.BindingTypeEnum(d.Get("binding_type").(string))
	r.BindingType = &bt

	nip := api.NameIdPolicyEnum(d.Get("name_id_policy").(string))
	r.NameIdPolicy.Set(&nip)

	da := api.DigestAlgorithmEnum(d.Get("digest_algorithm").(string))
	r.DigestAlgorithm = &da

	sa := api.SignatureAlgorithmEnum(d.Get("signature_algorithm").(string))
	r.SignatureAlgorithm = &sa

	if s, sok := d.GetOk("slo_url"); sok && s.(string) != "" {
		r.SloUrl.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("signing_kp"); sok && s.(string) != "" {
		r.SigningKp.Set(stringToPointer(s.(string)))
	}

	umm := api.UserMatchingModeEnum(d.Get("user_matching_mode").(string))
	r.UserMatchingMode.Set(&umm)
	return &r
}

func resourceSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlCreate(ctx).SAMLSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "uuid", res.Pk)

	if res.AuthenticationFlow.IsSet() {
		setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	}
	if res.EnrollmentFlow.IsSet() {
		setWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	}
	setWrapper(d, "enabled", res.Enabled)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "user_matching_mode", res.UserMatchingMode.Get())

	setWrapper(d, "pre_authentication_flow", res.PreAuthenticationFlow)
	setWrapper(d, "issuer", res.Issuer)
	setWrapper(d, "sso_url", res.SsoUrl)
	if res.SloUrl.IsSet() {
		setWrapper(d, "slo_url", res.SloUrl.Get())
	}
	setWrapper(d, "allow_idp_initiated", res.AllowIdpInitiated)
	setWrapper(d, "name_id_policy", res.NameIdPolicy.Get())
	setWrapper(d, "binding_type", res.BindingType)
	if res.SigningKp.IsSet() {
		setWrapper(d, "signing_kp", res.SigningKp.Get())
	}
	setWrapper(d, "digest_algorithm", res.DigestAlgorithm)
	setWrapper(d, "signature_algorithm", res.SignatureAlgorithm)
	setWrapper(d, "temporary_user_delete_after", res.TemporaryUserDeleteAfter)
	return diags
}

func resourceSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlUpdate(ctx, d.Id()).SAMLSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
