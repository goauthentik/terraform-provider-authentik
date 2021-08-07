package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	r.NameIdPolicy = &nip

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
	return &r
}

func resourceSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlCreate(ctx).SAMLSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	d.Set("uuid", res.Pk)

	if res.AuthenticationFlow.IsSet() {
		d.Set("authentication_flow", res.AuthenticationFlow.Get())
	}
	if res.EnrollmentFlow.IsSet() {
		d.Set("enrollment_flow", res.EnrollmentFlow.Get())
	}
	d.Set("enabled", res.Enabled)
	d.Set("policy_engine_mode", res.PolicyEngineMode)

	d.Set("pre_authentication_flow", res.PreAuthenticationFlow)
	d.Set("issuer", res.Issuer)
	d.Set("sso_url", res.SsoUrl)
	if res.SloUrl.IsSet() {
		d.Set("slo_url", res.SloUrl.Get())
	}
	d.Set("allow_idp_initiated", res.AllowIdpInitiated)
	d.Set("name_id_policy", res.NameIdPolicy)
	d.Set("binding_type", res.BindingType)
	if res.SigningKp.IsSet() {
		d.Set("signing_kp", res.SigningKp.Get())
	}
	d.Set("digest_algorithm", res.DigestAlgorithm)
	d.Set("signature_algorithm", res.SignatureAlgorithm)
	d.Set("temporary_user_delete_after", res.TemporaryUserDeleteAfter)
	return diags
}

func resourceSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceSAMLSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesSamlUpdate(ctx, d.Id()).SAMLSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceSAMLRead(ctx, d, m)
}

func resourceSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
