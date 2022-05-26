package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceProviderSAML() *schema.Resource {
	return &schema.Resource{
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
			"authorization_flow": {
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=3",
			},
			"assertion_valid_not_on_or_after": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=3",
			},
			"session_valid_not_on_or_after": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=3",
			},
			"name_id_mapping": {
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
			"signing_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verification_kp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sp_binding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.SPBINDINGENUM_REDIRECT,
			},
		},
	}
}

func resourceProviderSAMLSchemaToProvider(d *schema.ResourceData) *api.SAMLProviderRequest {
	r := api.SAMLProviderRequest{
		Name:                       d.Get("name").(string),
		AuthorizationFlow:          d.Get("authorization_flow").(string),
		AcsUrl:                     d.Get("acs_url").(string),
		Audience:                   stringToPointer(d.Get("audience").(string)),
		Issuer:                     stringToPointer(d.Get("issuer").(string)),
		AssertionValidNotBefore:    stringToPointer(d.Get("assertion_valid_not_before").(string)),
		AssertionValidNotOnOrAfter: stringToPointer(d.Get("assertion_valid_not_on_or_after").(string)),
		SessionValidNotOnOrAfter:   stringToPointer(d.Get("session_valid_not_on_or_after").(string)),
	}

	if s, sok := d.GetOk("name_id_mapping"); sok && s.(string) != "" {
		r.NameIdMapping.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("signing_kp"); sok && s.(string) != "" {
		r.SigningKp.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("verification_kp"); sok && s.(string) != "" {
		r.VerificationKp.Set(stringToPointer(s.(string)))
	}

	digA := d.Get("digest_algorithm").(string)
	a := api.DigestAlgorithmEnum(digA)
	r.DigestAlgorithm = &a

	sigA := d.Get("signature_algorithm").(string)
	c := api.SignatureAlgorithmEnum(sigA)
	r.SignatureAlgorithm = &c

	binding := d.Get("sp_binding").(string)
	j := api.SpBindingEnum(binding)
	r.SpBinding.Set(&j)

	r.PropertyMappings = sliceToString(d.Get("property_mappings").([]interface{}))

	return &r
}

func resourceProviderSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersSamlCreate(ctx).SAMLProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSAMLRead(ctx, d, m)
}

func resourceProviderSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersSamlRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("authorization_flow", res.AuthorizationFlow)
	localMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	d.Set("property_mappings", stringListConsistentMerge(localMappings, res.PropertyMappings))

	d.Set("acs_url", res.AcsUrl)
	d.Set("audience", res.Audience)
	d.Set("issuer", res.Issuer)
	d.Set("sp_binding", res.SpBinding.Get())
	d.Set("assertion_valid_not_before", res.AssertionValidNotBefore)
	d.Set("assertion_valid_not_on_or_after", res.AssertionValidNotOnOrAfter)
	d.Set("session_valid_not_on_or_after", res.SessionValidNotOnOrAfter)
	if res.NameIdMapping.IsSet() {
		d.Set("name_id_mapping", res.NameIdMapping.Get())
	}
	if res.SigningKp.IsSet() {
		d.Set("signing_kp", res.SigningKp.Get())
	}
	if res.VerificationKp.IsSet() {
		d.Set("verification_kp", res.VerificationKp.Get())
	}
	d.Set("digest_algorithm", res.DigestAlgorithm)
	d.Set("signature_algorithm", res.SignatureAlgorithm)

	return diags
}

func resourceProviderSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersSamlUpdate(ctx, int32(id)).SAMLProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSAMLRead(ctx, d, m)
}

func resourceProviderSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersSamlDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
