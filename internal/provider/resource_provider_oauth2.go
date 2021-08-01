package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"jwt_alg": {
				Type:     schema.TypeString,
				Default:  api.JWTALGENUM_HS256,
				Optional: true,
			},
			"rsa_key": {
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
		},
	}
}

func resourceProviderOAuth2SchemaToProvider(d *schema.ResourceData) (*api.OAuth2ProviderRequest, diag.Diagnostics) {
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

	subMode := d.Get("sub_mode").(string)
	a := api.SubModeEnum(subMode)
	r.SubMode = &a

	clientType := d.Get("client_type").(string)
	c := api.ClientTypeEnum(clientType)
	r.ClientType = &c

	jwtAlg := d.Get("jwt_alg").(string)
	j := api.JwtAlgEnum(jwtAlg)
	r.JwtAlg = &j

	redirect := make([]string, 0)
	for _, redirectS := range d.Get("redirect_uris").([]interface{}) {
		redirect = append(redirect, redirectS.(string))
	}
	r.RedirectUris = stringToPointer(strings.Join(redirect, "\n"))

	propertyMapping := make([]string, 0)
	for _, propertyMappingS := range d.Get("property_mappings").([]interface{}) {
		propertyMapping = append(propertyMapping, propertyMappingS.(string))
	}
	r.PropertyMappings = &propertyMapping

	return &r, nil
}

func resourceProviderOAuth2Create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	r, diags := resourceProviderOAuth2SchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Create(ctx).OAuth2ProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Retrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	d.Set("access_code_validity", res.AccessCodeValidity)
	d.Set("authorization_flow", res.AuthorizationFlow)
	d.Set("client_id", res.ClientId)
	d.Set("client_secret", res.ClientSecret)
	d.Set("client_type", res.ClientType)
	d.Set("include_claims_in_id_token", res.IncludeClaimsInIdToken)
	d.Set("issuer_mode", res.IssuerMode)
	d.Set("jwt_alg", res.JwtAlg)
	d.Set("property_mappings", res.PropertyMappings)
	if stringPointerResolve(res.RedirectUris) != "" {
		d.Set("redirect_uris", strings.Split(stringPointerResolve(res.RedirectUris), "\n"))
	} else {
		d.Set("redirect_uris", []string{})
	}
	if res.RsaKey.IsSet() {
		d.Set("rsa_key", res.RsaKey.Get())
	}
	d.Set("sub_mode", res.SubMode)
	d.Set("token_validity", res.TokenValidity)
	return diags
}

func resourceProviderOAuth2Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app, di := resourceProviderOAuth2SchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Update(ctx, int32(id)).OAuth2ProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersOauth2Destroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
