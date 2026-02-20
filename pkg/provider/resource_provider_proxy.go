package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderProxy() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
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
			"invalidation_flow": {
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
				DiffSuppressFunc: helpers.DiffSuppressExpression,
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.PROXYMODE_PROXY,
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedProxyModeEnumValues),
				Description:      helpers.EnumToDescription(api.AllowedProxyModeEnumValues),
			},
			"cookie_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token_validity": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=10",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"refresh_token_validity": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "days=30",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"jwks_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Deprecated. Use `jwt_federation_sources` instead.",
			},
			"jwt_federation_sources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "JWTs issued by keys configured in any of the selected sources can be used to authenticate on behalf of this provider.",
			},
			"jwt_federation_providers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "JWTs issued by any of the configured providers can be used to authenticate on behalf of this provider.",
			},
		},
	}
}

func resourceProviderProxySchemaToProvider(d *schema.ResourceData) *api.ProxyProviderRequest {
	r := api.ProxyProviderRequest{
		Name:                      d.Get("name").(string),
		AuthorizationFlow:         d.Get("authorization_flow").(string),
		InvalidationFlow:          d.Get("invalidation_flow").(string),
		ExternalHost:              d.Get("external_host").(string),
		Mode:                      api.ProxyMode(d.Get("mode").(string)).Ptr(),
		PropertyMappings:          helpers.CastSlice[string](d, "property_mappings"),
		JwtFederationSources:      helpers.CastSlice[string](d, "jwt_federation_sources"),
		AuthenticationFlow:        *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		InternalHost:              helpers.GetP[string](d, "internal_host"),
		InternalHostSslValidation: helpers.GetP[bool](d, "internal_host_ssl_validation"),

		SkipPathRegex: helpers.GetP[string](d, "skip_path_regex"),

		BasicAuthEnabled:           helpers.GetP[bool](d, "basic_auth_enabled"),
		InterceptHeaderAuth:        helpers.GetP[bool](d, "intercept_header_auth"),
		BasicAuthUserAttribute:     helpers.GetP[string](d, "basic_auth_username_attribute"),
		BasicAuthPasswordAttribute: helpers.GetP[string](d, "basic_auth_password_attribute"),

		CookieDomain: helpers.GetP[string](d, "cookie_domain"),

		AccessTokenValidity:  helpers.GetP[string](d, "access_token_validity"),
		RefreshTokenValidity: helpers.GetP[string](d, "refresh_token_validity"),
	}

	providers := d.Get("jwt_federation_providers").([]any)
	r.JwtFederationProviders = make([]int32, len(providers))
	for i, prov := range providers {
		r.JwtFederationProviders[i] = int32(prov.(int))
	}

	return &r
}

func resourceProviderProxyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderProxySchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersProxyCreate(ctx).ProxyProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersProxyRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "client_id", res.ClientId)
	helpers.SetWrapper(d, "intercept_header_auth", res.InterceptHeaderAuth)
	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "invalidation_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "internal_host", res.InternalHost)
	helpers.SetWrapper(d, "external_host", res.ExternalHost)
	helpers.SetWrapper(d, "internal_host_ssl_validation", res.InternalHostSslValidation)
	helpers.SetWrapper(d, "skip_path_regex", res.SkipPathRegex)
	helpers.SetWrapper(d, "basic_auth_enabled", res.BasicAuthEnabled)
	helpers.SetWrapper(d, "basic_auth_username_attribute", res.BasicAuthUserAttribute)
	helpers.SetWrapper(d, "basic_auth_password_attribute", res.BasicAuthPasswordAttribute)
	helpers.SetWrapper(d, "mode", res.Mode)
	helpers.SetWrapper(d, "cookie_domain", res.CookieDomain)
	helpers.SetWrapper(d, "access_token_validity", res.AccessTokenValidity)
	helpers.SetWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	localMappings := helpers.CastSlice[string](d, "property_mappings")
	if len(localMappings) > 0 {
		// Only update mappings if any were set in TF resource, since authentik will always set the
		// default mappings, even when nothing is specified
		helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
			localMappings,
			res.PropertyMappings,
		))
	}
	helpers.SetWrapper(d, "jwt_federation_providers", helpers.ListConsistentMerge(
		helpers.CastSlice[int](d, "jwt_federation_providers"),
		helpers.Slice32ToInt(res.JwtFederationProviders),
	))
	helpers.SetWrapper(d, "jwt_federation_sources", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "jwt_federation_sources"),
		res.JwtFederationSources,
	))
	return diags
}

func resourceProviderProxyUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderProxySchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersProxyUpdate(ctx, int32(id)).ProxyProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderProxyRead(ctx, d, m)
}

func resourceProviderProxyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersProxyDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
