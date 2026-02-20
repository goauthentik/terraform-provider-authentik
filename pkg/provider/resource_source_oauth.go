package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
			"promoted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"authorization_code_auth_method": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHORIZATIONCODEAUTHMETHODENUM_BASIC_AUTH,
				Description:      helpers.EnumToDescription(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedAuthorizationCodeAuthMethodEnumEnumValues),
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"user_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserMatchingModeEnumEnumValues),
			},
			"group_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.GROUPMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedGroupMatchingModeEnumEnumValues),
			},

			"provider_type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedProviderTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedProviderTypeEnumEnumValues),
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
				Description:      "Manually configure JWKS keys for use with machine-to-machine authentication. " + helpers.JSONDescription,
				Computed:         true,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"pkce": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.PKCEMETHODENUM_NONE,
				Description:      helpers.EnumToDescription(api.AllowedPKCEMethodEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPKCEMethodEnumEnumValues),
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
		Enabled:          new(d.Get("enabled").(bool)),
		Promoted:         new(d.Get("promoted").(bool)),
		UserPathTemplate: new(d.Get("user_path_template").(string)),

		ProviderType:                api.ProviderTypeEnum(d.Get("provider_type").(string)),
		ConsumerKey:                 d.Get("consumer_key").(string),
		ConsumerSecret:              d.Get("consumer_secret").(string),
		AuthorizationCodeAuthMethod: api.AuthorizationCodeAuthMethodEnum(d.Get("authorization_code_auth_method").(string)).Ptr(),
		PolicyEngineMode:            api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:            api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode:           api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),
		AuthenticationFlow:          *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		EnrollmentFlow:              *api.NewNullableString(helpers.GetP[string](d, "enrollment_flow")),

		RequestTokenUrl:       *api.NewNullableString(helpers.GetP[string](d, "request_token_url")),
		AuthorizationUrl:      *api.NewNullableString(helpers.GetP[string](d, "authorization_url")),
		AccessTokenUrl:        *api.NewNullableString(helpers.GetP[string](d, "access_token_url")),
		ProfileUrl:            *api.NewNullableString(helpers.GetP[string](d, "profile_url")),
		AdditionalScopes:      helpers.GetP[string](d, "additional_scopes"),
		OidcWellKnownUrl:      helpers.GetP[string](d, "oidc_well_known_url"),
		OidcJwksUrl:           helpers.GetP[string](d, "oidc_jwks_url"),
		Pkce:                  helpers.CastString[api.PKCEMethodEnum](helpers.GetP[string](d, "pkce")),
		UserPropertyMappings:  helpers.CastSlice[string](d, "property_mappings"),
		GroupPropertyMappings: helpers.CastSlice[string](d, "property_mappings_group"),
	}

	jwks, err := helpers.GetJSON[map[string]any](d, ("oidc_jwks"))
	r.OidcJwks = jwks
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func resourceSourceOAuthCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceSourceOAuthSchemaToSource(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.SourcesApi.SourcesOauthCreate(ctx).OAuthSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceOAuthRead(ctx, d, m)
}

func resourceSourceOAuthRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesOauthRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "user_path_template", res.UserPathTemplate)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "promoted", res.Promoted)

	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	helpers.SetWrapper(d, "authorization_code_auth_method", res.AuthorizationCodeAuthMethod)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "user_matching_mode", res.UserMatchingMode)
	helpers.SetWrapper(d, "group_matching_mode", res.GroupMatchingMode)
	helpers.SetWrapper(d, "additional_scopes", res.AdditionalScopes)
	helpers.SetWrapper(d, "provider_type", res.ProviderType)
	helpers.SetWrapper(d, "consumer_key", res.ConsumerKey)
	helpers.SetWrapper(d, "request_token_url", res.RequestTokenUrl.Get())
	helpers.SetWrapper(d, "authorization_url", res.AuthorizationUrl.Get())
	helpers.SetWrapper(d, "access_token_url", res.AccessTokenUrl.Get())
	helpers.SetWrapper(d, "profile_url", res.ProfileUrl.Get())
	helpers.SetWrapper(d, "callback_uri", res.CallbackUrl)
	helpers.SetWrapper(d, "oidc_well_known_url", res.GetOidcWellKnownUrl())
	helpers.SetWrapper(d, "oidc_jwks_url", res.GetOidcJwksUrl())
	helpers.SetWrapper(d, "pkce", res.GetPkce())
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.UserPropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.GroupPropertyMappings,
	))
	return helpers.SetJSON(d, "oidc_jwks", res.OidcJwks)
}

func resourceSourceOAuthUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app, diags := resourceSourceOAuthSchemaToSource(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.SourcesApi.SourcesOauthUpdate(ctx, d.Id()).OAuthSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceOAuthRead(ctx, d, m)
}

func resourceSourceOAuthDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesOauthDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
