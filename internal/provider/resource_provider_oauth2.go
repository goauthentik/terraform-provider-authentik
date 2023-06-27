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
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
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
				Computed:  true,
			},
			"access_code_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=1",
			},
			"access_token_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=10",
			},
			"refresh_token_validity": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=30",
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
			"jwks_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "JWTs issued by keys configured in any of the selected sources can be used to authenticate on behalf of this provider.",
			},
		},
	}
}

func resourceProviderOAuth2SchemaToProvider(d *schema.ResourceData) *api.OAuth2ProviderRequest {
	r := api.OAuth2ProviderRequest{
		Name:                   d.Get("name").(string),
		AuthorizationFlow:      d.Get("authorization_flow").(string),
		AccessCodeValidity:     api.PtrString(d.Get("access_code_validity").(string)),
		AccessTokenValidity:    api.PtrString(d.Get("access_token_validity").(string)),
		RefreshTokenValidity:   api.PtrString(d.Get("refresh_token_validity").(string)),
		IncludeClaimsInIdToken: api.PtrBool(d.Get("include_claims_in_id_token").(bool)),
		ClientId:               api.PtrString(d.Get("client_id").(string)),
		IssuerMode:             api.IssuerModeEnum(d.Get("issuer_mode").(string)).Ptr(),
		SubMode:                api.SubModeEnum(d.Get("sub_mode").(string)).Ptr(),
		ClientType:             api.ClientTypeEnum(d.Get("client_type").(string)).Ptr(),
		PropertyMappings:       castSlice[string](d.Get("property_mappings").([]interface{})),
		JwksSources:            castSlice[string](d.Get("jwks_sources").([]interface{})),
	}

	if s, sok := d.GetOk("authentication_flow"); sok && s.(string) != "" {
		r.AuthenticationFlow.Set(api.PtrString(s.(string)))
	}
	if s, sok := d.GetOk("client_secret"); sok && s.(string) != "" {
		r.ClientSecret = api.PtrString(s.(string))
	}

	if s, sok := d.GetOk("signing_key"); sok && s.(string) != "" {
		r.SigningKey.Set(api.PtrString(s.(string)))
	}

	redirectUris := castSlice[string](d.Get("redirect_uris").([]interface{}))
	r.RedirectUris = api.PtrString(strings.Join(redirectUris, "\n"))
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

	setWrapper(d, "name", res.Name)
	setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	setWrapper(d, "authorization_flow", res.AuthorizationFlow)
	setWrapper(d, "client_id", res.ClientId)
	setWrapper(d, "client_secret", res.ClientSecret)
	setWrapper(d, "client_type", res.ClientType)
	setWrapper(d, "include_claims_in_id_token", res.IncludeClaimsInIdToken)
	setWrapper(d, "issuer_mode", res.IssuerMode)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.PropertyMappings))
	if *res.RedirectUris != "" {
		setWrapper(d, "redirect_uris", strings.Split(*res.RedirectUris, "\n"))
	} else {
		setWrapper(d, "redirect_uris", []string{})
	}
	if res.SigningKey.IsSet() {
		setWrapper(d, "signing_key", res.SigningKey.Get())
	}
	setWrapper(d, "sub_mode", res.SubMode)
	setWrapper(d, "access_code_validity", res.AccessCodeValidity)
	setWrapper(d, "access_token_validity", res.AccessTokenValidity)
	setWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	localJWKSSources := castSlice[string](d.Get("jwks_sources").([]interface{}))
	setWrapper(d, "jwks_sources", listConsistentMerge(localJWKSSources, res.JwksSources))
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
