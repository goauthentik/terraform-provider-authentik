package authentik

import (
	"context"
	"strconv"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"skip_path_regex": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceProviderProxySchemaToProvider(d *schema.ResourceData) (*api.ProxyProviderRequest, diag.Diagnostics) {
	r := api.ProxyProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("authorization_flow").(string),
		ExternalHost:      d.Get("external_host").(string),
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
	if l, ok := d.Get("basic_auth_username_attribute").(string); ok {
		r.BasicAuthUserAttribute = &l
	}
	if l, ok := d.Get("basic_auth_password_attribute").(string); ok {
		r.BasicAuthPasswordAttribute = &l
	}

	if l, ok := d.Get("cookie_domain").(string); ok {
		r.CookieDomain = &l
	}

	pm := d.Get("mode").(string)
	var pma api.ProxyMode
	switch pm {
	case string(api.PROXYMODE_PROXY):
		pma = api.PROXYMODE_PROXY
	case string(api.PROXYMODE_FORWARD_SINGLE):
		pma = api.PROXYMODE_FORWARD_SINGLE
	case string(api.PROXYMODE_FORWARD_DOMAIN):
		pma = api.PROXYMODE_FORWARD_DOMAIN
	default:
		return nil, diag.Errorf("invalid mode %s", pm)
	}
	r.Mode = &pma
	return &r, nil
}

func resourceProviderProxyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	r, diags := resourceProviderProxySchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersProxyCreate(ctx).ProxyProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersProxyRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	d.Set("authorization_flow", res.AuthorizationFlow)
	d.Set("internal_host", res.InternalHost)
	d.Set("external_host", res.ExternalHost)
	d.Set("internal_host_ssl_validation", res.InternalHostSslValidation)
	d.Set("skip_path_regex", res.SkipPathRegex)
	d.Set("basic_auth_enabled", res.BasicAuthEnabled)
	d.Set("basic_auth_username_attribute", res.BasicAuthUserAttribute)
	d.Set("basic_auth_password_attribute", res.BasicAuthPasswordAttribute)
	d.Set("mode", res.Mode)
	d.Set("cookie_domain", res.CookieDomain)
	return diags
}

func resourceProviderProxyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app, di := resourceProviderProxySchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.ProvidersApi.ProvidersProxyUpdate(ctx, int32(id)).ProxyProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersProxyDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
