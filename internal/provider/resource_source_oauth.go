package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceOAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceOAuthCreate,
		ReadContext:   resourceSourceOAuthRead,
		UpdateContext: resourceSourceOAuthUpdate,
		DeleteContext: resourceSourceOAuthDelete,
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

			"provider_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_token_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_token_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"additional_scopes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consumer_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumer_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"callback_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSourceOAuthSchemaToSource(d *schema.ResourceData) *api.OAuthSourceRequest {
	r := api.OAuthSourceRequest{
		Name:    d.Get("name").(string),
		Slug:    d.Get("slug").(string),
		Enabled: boolToPointer(d.Get("enabled").(bool)),

		ProviderType:   api.ProviderTypeEnum(d.Get("provider_type").(string)),
		ConsumerKey:    d.Get("consumer_key").(string),
		ConsumerSecret: d.Get("consumer_secret").(string),
	}

	r.AuthenticationFlow.Set(stringToPointer(d.Get("authentication_flow").(string)))
	r.EnrollmentFlow.Set(stringToPointer(d.Get("enrollment_flow").(string)))

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	r.PolicyEngineMode = &pm

	umm := api.UserMatchingModeEnum(d.Get("user_matching_mode").(string))
	r.UserMatchingMode = &umm

	if s, sok := d.GetOk("request_token_url"); sok && s.(string) != "" {
		r.RequestTokenUrl.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("authorization_url"); sok && s.(string) != "" {
		r.AuthorizationUrl.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("access_token_url"); sok && s.(string) != "" {
		r.AccessTokenUrl.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("profile_url"); sok && s.(string) != "" {
		r.ProfileUrl.Set(stringToPointer(s.(string)))
	}
	if s, sok := d.GetOk("additional_scopes"); sok && s.(string) != "" {
		r.AdditionalScopes = stringToPointer(s.(string))
	}
	return &r
}

func resourceSourceOAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceOAuthSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesOauthCreate(ctx).OAuthSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceOAuthRead(ctx, d, m)
}

func resourceSourceOAuthRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesOauthRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
	d.Set("user_matching_mode", res.UserMatchingMode)
	d.Set("additional_scopes", res.AdditionalScopes)
	d.Set("provider_type", res.ProviderType)
	d.Set("consumer_key", res.ConsumerKey)
	if res.RequestTokenUrl.IsSet() {
		d.Set("request_token_url", res.RequestTokenUrl.Get())
	}
	if res.AuthorizationUrl.IsSet() {
		d.Set("authorization_url", res.AuthorizationUrl.Get())
	}
	if res.AccessTokenUrl.IsSet() {
		d.Set("access_token_url", res.AccessTokenUrl.Get())
	}
	if res.ProfileUrl.IsSet() {
		d.Set("profile_url", res.ProfileUrl.Get())
	}
	d.Set("callback_uri", res.CallbackUrl)
	return diags
}

func resourceSourceOAuthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceOAuthSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesOauthUpdate(ctx, d.Id()).OAuthSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceOAuthRead(ctx, d, m)
}

func resourceSourceOAuthDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesOauthDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
