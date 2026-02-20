package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderOAuth2() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
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
			"invalidation_flow": {
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CLIENTTYPEENUM_CONFIDENTIAL,
				Description:      helpers.EnumToDescription(api.AllowedClientTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedClientTypeEnumEnumValues),
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=1",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
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
			"refresh_token_threshold": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "seconds=0",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
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
			"encryption_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allowed_redirect_uris": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"logout_method": {
				Type:             schema.TypeString,
				Default:          api.OAUTH2PROVIDERLOGOUTMETHODENUM_BACKCHANNEL,
				Optional:         true,
				Description:      helpers.EnumToDescription(api.AllowedOAuth2ProviderLogoutMethodEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedOAuth2ProviderLogoutMethodEnumEnumValues),
			},
			"logout_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_mode": {
				Type:             schema.TypeString,
				Default:          api.SUBMODEENUM_HASHED_USER_ID,
				Optional:         true,
				Description:      helpers.EnumToDescription(api.AllowedSubModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSubModeEnumEnumValues),
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Default:          api.ISSUERMODEENUM_PER_PROVIDER,
				Optional:         true,
				Description:      helpers.EnumToDescription(api.AllowedIssuerModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedIssuerModeEnumEnumValues),
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

func resourceProviderOAuth2SchemaToProvider(d *schema.ResourceData) *api.OAuth2ProviderRequest {
	r := api.OAuth2ProviderRequest{
		Name:                   d.Get("name").(string),
		AuthorizationFlow:      d.Get("authorization_flow").(string),
		AuthenticationFlow:     *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		InvalidationFlow:       d.Get("invalidation_flow").(string),
		AccessCodeValidity:     new(d.Get("access_code_validity").(string)),
		AccessTokenValidity:    new(d.Get("access_token_validity").(string)),
		RefreshTokenValidity:   new(d.Get("refresh_token_validity").(string)),
		RefreshTokenThreshold:  helpers.GetP[string](d, "refresh_token_threshold"),
		IncludeClaimsInIdToken: new(d.Get("include_claims_in_id_token").(bool)),
		ClientId:               new(d.Get("client_id").(string)),
		ClientSecret:           helpers.GetP[string](d, "client_secret"),
		IssuerMode:             api.IssuerModeEnum(d.Get("issuer_mode").(string)).Ptr(),
		SubMode:                api.SubModeEnum(d.Get("sub_mode").(string)).Ptr(),
		ClientType:             api.ClientTypeEnum(d.Get("client_type").(string)).Ptr(),
		PropertyMappings:       helpers.CastSlice[string](d, "property_mappings"),
		JwtFederationSources:   helpers.CastSlice[string](d, "jwt_federation_sources"),
		LogoutMethod:           helpers.CastString[api.OAuth2ProviderLogoutMethodEnum](helpers.GetP[string](d, "logout_method")),
		LogoutUri:              helpers.GetP[string](d, "logout_uri"),

		SigningKey:             *api.NewNullableString(helpers.GetP[string](d, "signing_key")),
		EncryptionKey:          *api.NewNullableString(helpers.GetP[string](d, "encryption_key")),
		RedirectUris:           listToRedirectURIsRequest(d.Get("allowed_redirect_uris").([]any)),
		JwtFederationProviders: helpers.CastSliceInt32(d.Get("jwt_federation_providers").([]any)),
	}
	return &r
}

func listToRedirectURIsRequest(raw []any) []api.RedirectURIRequest {
	rus := []api.RedirectURIRequest{}
	for _, rr := range raw {
		rd := rr.(map[string]any)
		rus = append(rus, api.RedirectURIRequest{
			MatchingMode: api.MatchingModeEnum(rd["matching_mode"].(string)),
			Url:          rd["url"].(string),
		})
	}
	return rus
}

func listToRedirectURIs(raw []any) []api.RedirectURI {
	rus := []api.RedirectURI{}
	for _, rr := range raw {
		rd := rr.(map[string]any)
		rus = append(rus, api.RedirectURI{
			MatchingMode: api.MatchingModeEnum(rd["matching_mode"].(string)),
			Url:          rd["url"].(string),
		})
	}
	return rus
}

func redirectURIsToList(raw []api.RedirectURI) []map[string]any {
	rus := []map[string]any{}
	for _, rr := range raw {
		rus = append(rus, map[string]any{
			"matching_mode": string(rr.MatchingMode),
			"url":           rr.Url,
		})
	}
	return rus
}

func resourceProviderOAuth2Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Create(ctx).OAuth2ProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Retrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "invalidation_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "client_id", res.ClientId)
	helpers.SetWrapper(d, "client_secret", res.ClientSecret)
	helpers.SetWrapper(d, "client_type", res.ClientType)
	helpers.SetWrapper(d, "include_claims_in_id_token", res.IncludeClaimsInIdToken)
	helpers.SetWrapper(d, "issuer_mode", res.IssuerMode)
	helpers.SetWrapper(d, "logout_method", res.LogoutMethod)
	helpers.SetWrapper(d, "logout_uri", res.LogoutUri)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	helpers.SetWrapper(d, "allowed_redirect_uris", redirectURIsToList(
		helpers.ListConsistentMerge(
			listToRedirectURIs(d.Get("allowed_redirect_uris").([]any)),
			res.RedirectUris,
		),
	))
	helpers.SetWrapper(d, "signing_key", res.SigningKey.Get())
	helpers.SetWrapper(d, "encryption_key", res.EncryptionKey.Get())
	helpers.SetWrapper(d, "sub_mode", res.SubMode)
	helpers.SetWrapper(d, "access_code_validity", res.AccessCodeValidity)
	helpers.SetWrapper(d, "access_token_validity", res.AccessTokenValidity)
	helpers.SetWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	helpers.SetWrapper(d, "refresh_token_threshold", res.RefreshTokenThreshold)
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

func resourceProviderOAuth2Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderOAuth2SchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersOauth2Update(ctx, int32(id)).OAuth2ProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderOAuth2Read(ctx, d, m)
}

func resourceProviderOAuth2Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersOauth2Destroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
