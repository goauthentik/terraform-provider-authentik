package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceProviderOAuth2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProviderOAuth2Create,
		ReadContext:   resourceProviderOAuth2Read,
		UpdateContext: resourceProviderOAuth2Update,
		DeleteContext: resourceProviderOAuth2Delete,
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
			"client_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.CLIENTTYPEENUM_CONFIDENTIAL,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"access_code_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=1",
			},
			"token_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=10",
			},
			"include_claims_in_id_token": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"signing_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"verification_keys": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sub_mode": {
				Type:     schema.TypeString,
				Default:  api.SUBMODEENUM_HASHED_USER_ID,
				Optional: true,
			},
			"issuer_mode": {
				Type:     schema.TypeString,
				Default:  api.ISSUERMODEENUM_PER_PROVIDER,
				Optional: true,
			},
		},
	}
}

func resourceProviderOAuth2SchemaToProvider(d *schema.ResourceData) *api.OAuth2ProviderRequest {
	r := api.OAuth2ProviderRequest{
		Name:                   d.Get("name").(string),
		AuthorizationFlow:      d.Get("authorization_flow").(string),
		AccessCodeValidity:     stringToPointer(d.Get("access_code_validity").(string)),
		TokenValidity:          stringToPointer(d.Get("token_validity").(string)),
		IncludeClaimsInIdToken: boolToPointer(d.Get("include_claims_in_id_token").(bool)),
		ClientId:               stringToPointer(d.Get("client_id").(string)),
	}

	if s, sok := d.GetOk("client_secret"); sok && s.(string) != "" {
		r.ClientSecret = stringToPointer(s.(string))
	}

	if s, sok := d.GetOk("signing_key"); sok && s.(string) != "" {
		r.SigningKey.Set(stringToPointer(s.(string)))
	}

	subMode := d.Get("sub_mode").(string)
	a := api.SubModeEnum(subMode)
	r.SubMode.Set(&a)

	clientType := d.Get("client_type").(string)
	c := api.ClientTypeEnum(clientType)
	r.ClientType.Set(&c)

	redirectUris := sliceToString(d.Get("redirect_uris").([]interface{}))
	r.RedirectUris = stringToPointer(strings.Join(redirectUris, "\n"))

	if vk, sok := d.GetOk("verification_keys"); sok {
		r.SetVerificationKeys(sliceToString(vk.([]interface{})))
	}

	r.PropertyMappings = sliceToString(d.Get("property_mappings").([]interface{}))
	return &r
}

func resourceProviderOAuth2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Create(ctx).OAuth2ProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Retrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("access_code_validity", res.AccessCodeValidity)
	d.Set("authorization_flow", res.AuthorizationFlow)
	d.Set("client_id", res.ClientId)
	d.Set("client_secret", res.ClientSecret)
	d.Set("client_type", res.ClientType.Get())
	d.Set("include_claims_in_id_token", res.IncludeClaimsInIdToken)
	d.Set("issuer_mode", res.IssuerMode.Get())
	localMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	d.Set("property_mappings", stringListConsistentMerge(localMappings, res.PropertyMappings))
	if stringPointerResolve(res.RedirectUris) != "" {
		d.Set("redirect_uris", strings.Split(stringPointerResolve(res.RedirectUris), "\n"))
	} else {
		d.Set("redirect_uris", []string{})
	}
	d.Set("verification_keys", res.GetVerificationKeys())
	if res.SigningKey.IsSet() {
		d.Set("signing_key", res.SigningKey.Get())
	}
	d.Set("sub_mode", res.SubMode.Get())
	d.Set("token_validity", res.TokenValidity)
	return diags
}

func resourceProviderOAuth2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Update(ctx, int32(id)).OAuth2ProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersOauth2Destroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
