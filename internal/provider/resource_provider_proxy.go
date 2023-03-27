package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceProviderProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProviderProxyCreate,
		ReadContext:   resourceProviderProxyRead,
		UpdateContext: resourceProviderProxyUpdate,
		DeleteContext: resourceProviderProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_host_ssl_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"skip_path_regex": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
			"intercept_header_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"basic_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"basic_auth_username_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"basic_auth_password_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.PROXYMODE_PROXY,
			},
			"cookie_domain": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceProviderProxySchemaToProvider(d *schema.ResourceData) *api.ProxyProviderRequest {
	r := api.ProxyProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("authorization_flow").(string),
		ExternalHost:      d.Get("external_host").(string),
	}

	if s, sok := d.GetOk("authentication_flow"); sok && s.(string) != "" {
		r.AuthenticationFlow.Set(stringToPointer(s.(string)))
	}
	if l, ok := d.Get("internal_host").(string); ok {
		r.InternalHost = &l
	}
	if l, ok := d.Get("internal_host_ssl_validation").(bool); ok {
		r.InternalHostSslValidation = &l
	}

	if l, ok := d.Get("skip_path_regex").(string); ok {
		r.SkipPathRegex = &l
	}

	if l, ok := d.Get("basic_auth_enabled").(bool); ok {
		r.BasicAuthEnabled = &l
	}
	if l, ok := d.Get("intercept_header_auth").(bool); ok {
		r.InterceptHeaderAuth = &l
	}
	if l, ok := d.Get("basic_auth_username_attribute").(string); ok {
		r.BasicAuthUserAttribute = &l
	}
	if l, ok := d.Get("basic_auth_password_attribute").(string); ok {
		r.BasicAuthPasswordAttribute = &l
	}

	if l, ok := d.Get("cookie_domain").(string); ok {
		r.CookieDomain = &l
	}

	if l, ok := d.Get("access_token_validity").(string); ok {
		r.AccessTokenValidity = &l
	}
	if l, ok := d.Get("refresh_token_validity").(string); ok {
		r.RefreshTokenValidity = &l
	}

	pm := api.ProxyMode(d.Get("mode").(string))
	r.Mode.Set(&pm)

	r.PropertyMappings = sliceToString(d.Get("property_mappings").([]interface{}))
	r.JwksSources = sliceToString(d.Get("jwks_sources").([]interface{}))
	return &r
}

func resourceProviderProxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderProxySchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersProxyCreate(ctx).ProxyProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersProxyRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "client_id", res.ClientId)
	setWrapper(d, "intercept_header_auth", res.InterceptHeaderAuth)
	setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	setWrapper(d, "authorization_flow", res.AuthorizationFlow)
	setWrapper(d, "internal_host", res.InternalHost)
	setWrapper(d, "external_host", res.ExternalHost)
	setWrapper(d, "internal_host_ssl_validation", res.InternalHostSslValidation)
	setWrapper(d, "skip_path_regex", res.SkipPathRegex)
	setWrapper(d, "basic_auth_enabled", res.BasicAuthEnabled)
	setWrapper(d, "basic_auth_username_attribute", res.BasicAuthUserAttribute)
	setWrapper(d, "basic_auth_password_attribute", res.BasicAuthPasswordAttribute)
	setWrapper(d, "mode", res.Mode.Get())
	setWrapper(d, "cookie_domain", res.CookieDomain)
	setWrapper(d, "access_token_validity", res.AccessTokenValidity)
	setWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	localMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	if len(localMappings) > 0 {
		setWrapper(d, "property_mappings", stringListConsistentMerge(localMappings, res.PropertyMappings))
	}
	localJWKSSources := sliceToString(d.Get("jwks_sources").([]interface{}))
	setWrapper(d, "jwks_sources", stringListConsistentMerge(localJWKSSources, res.JwksSources))
	return diags
}

func resourceProviderProxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderProxySchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersProxyUpdate(ctx, int32(id)).ProxyProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersProxyDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
