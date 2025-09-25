package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceOAuth() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
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
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "goauthentik.io/sources/%(slug)s",
				Optional: true,
			},
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"authorization_code_auth_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHORIZATIONCODEAUTHMETHODENUM_BASIC_AUTH,
				Description:      EnumToDescription(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues),
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"user_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERMATCHINGMODEENUM_IDENTIFIER,
				Description:      EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserMatchingModeEnumEnumValues),
			},
			"group_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.GROUPMATCHINGMODEENUM_IDENTIFIER,
				Description:      EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedGroupMatchingModeEnumEnumValues),
			},

			"provider_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      EnumToDescription(api.AllowedProviderTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedProviderTypeEnumEnumValues),
			},

			"request_token_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manually configure OAuth2 URLs when `oidc_well_known_url` is not set.",
			},
			"authorization_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manually configure OAuth2 URLs when `oidc_well_known_url` is not set.",
			},
			"access_token_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Only required for OAuth1.",
			},
			"profile_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Manually configure OAuth2 URLs when `oidc_well_known_url` is not set.",
			},

			"oidc_well_known_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Automatically configure source from OIDC well-known endpoint. URL is taken as is, and should end with `.well-known/openid-configuration`.",
			},
			"oidc_jwks_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Automatically configure JWKS if not specified by `oidc_well_known_url`.",
			},
			"oidc_jwks": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Manually configure JWKS keys for use with machine-to-machine authentication. JSON format expected. Use jsonencode() to pass objects.",
				Computed:         true,
				DiffSuppressFunc: diffSuppressJSON,
				ValidateDiagFunc: ValidateJSON,
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
			"property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property_mappings_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSourceOAuthSchemaToSource(d *schema.ResourceData) (*api.OAuthSourceRequest, diag.Diagnostics) {
	r := api.OAuthSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate: api.PtrString(d.Get("user_path_template").(string)),

		ProviderType:                api.ProviderTypeEnum(d.Get("provider_type").(string)),
		ConsumerKey:                 d.Get("consumer_key").(string),
		ConsumerSecret:              d.Get("consumer_secret").(string),
		AuthorizationCodeAuthMethod: api.AuthorizationCodeAuthMethodEnum(d.Get("authorization_code_auth_method").(string)).Ptr(),
		PolicyEngineMode:            api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:            api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode:           api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),
		AuthenticationFlow:          *api.NewNullableString(getP[string](d, "authentication_flow")),
		EnrollmentFlow:              *api.NewNullableString(getP[string](d, "enrollment_flow")),

		RequestTokenUrl:  *api.NewNullableString(getP[string](d, "request_token_url")),
		AuthorizationUrl: *api.NewNullableString(getP[string](d, "authorization_url")),
		AccessTokenUrl:   *api.NewNullableString(getP[string](d, "access_token_url")),
		ProfileUrl:       *api.NewNullableString(getP[string](d, "profile_url")),
		AdditionalScopes: getP[string](d, "additional_scopes"),
		OidcWellKnownUrl: getP[string](d, "oidc_well_known_url"),
		OidcJwksUrl:      getP[string](d, "oidc_jwks_url"),

		UserPropertyMappings:  castSlice[string](d.Get("property_mappings").([]interface{})),
		GroupPropertyMappings: castSlice[string](d.Get("property_mappings_group").([]interface{})),
	}

	jwks, err := getJSON[map[string]interface{}](d, ("oidc_jwks"))
	r.OidcJwks = jwks
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func resourceSourceOAuthCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceSourceOAuthSchemaToSource(d)
	if diags != nil {
		return diags
	}

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

	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "user_path_template", res.UserPathTemplate)

	if res.AuthenticationFlow.IsSet() {
		setWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	}
	if res.EnrollmentFlow.IsSet() {
		setWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	}
	setWrapper(d, "enabled", res.Enabled)
	setWrapper(d, "authorization_code_auth_method", res.AuthorizationCodeAuthMethod)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "user_matching_mode", res.UserMatchingMode)
	setWrapper(d, "group_matching_mode", res.GroupMatchingMode)
	setWrapper(d, "additional_scopes", res.AdditionalScopes)
	setWrapper(d, "provider_type", res.ProviderType)
	setWrapper(d, "consumer_key", res.ConsumerKey)
	if res.RequestTokenUrl.IsSet() {
		setWrapper(d, "request_token_url", res.RequestTokenUrl.Get())
	}
	if res.AuthorizationUrl.IsSet() {
		setWrapper(d, "authorization_url", res.AuthorizationUrl.Get())
	}
	if res.AccessTokenUrl.IsSet() {
		setWrapper(d, "access_token_url", res.AccessTokenUrl.Get())
	}
	if res.ProfileUrl.IsSet() {
		setWrapper(d, "profile_url", res.ProfileUrl.Get())
	}
	setWrapper(d, "callback_uri", res.CallbackUrl)
	setWrapper(d, "oidc_well_known_url", res.GetOidcWellKnownUrl())
	setWrapper(d, "oidc_jwks_url", res.GetOidcJwksUrl())
	b, err := json.Marshal(res.GetOidcJwks())
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "oidc_jwks", string(b))
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.UserPropertyMappings))
	localGroupMappings := castSlice[string](d.Get("property_mappings_group").([]interface{}))
	setWrapper(d, "property_mappings_group", listConsistentMerge(localGroupMappings, res.GroupPropertyMappings))
	return diags
}

func resourceSourceOAuthUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app, diags := resourceSourceOAuthSchemaToSource(d)
	if diags != nil {
		return diags
	}

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
