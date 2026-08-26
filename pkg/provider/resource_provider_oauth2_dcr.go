package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderOAuth2DCR() *schema.Resource {
	return &schema.Resource{
		Description:   "Providers --- ",
		CreateContext: resourceProviderOAuth2DCRCreate,
		ReadContext:   resourceProviderOAuth2DCRRead,
		UpdateContext: resourceProviderOAuth2DCRUpdate,
		DeleteContext: resourceProviderOAuth2DCRDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"oauth2_provider": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "PK of the OAuth2 provider dynamic client registration applies to.",
			},
			"default_application_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group to assign to automatically created applications.",
			},
			"override_authorization_flow": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authorization flow applied to dynamically registered clients.",
			},
			"override_invalidation_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"override_property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"access_token_validity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Maximum access token validity for registered clients. Format: hours=1;minutes=2;seconds=3.",
			},
			"refresh_token_validity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Maximum refresh token validity for registered clients. Format: hours=1;minutes=2;seconds=3.",
			},
			"allowed_grant_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If empty, all grant types are allowed. " + helpers.EnumToDescription(api.AllowedGrantTypeEnumEnumValues),
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: helpers.StringInEnum(api.AllowedGrantTypeEnumEnumValues),
				},
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
		},
	}
}

func resourceProviderOAuth2DCRSchemaToModel(d *schema.ResourceData) *api.OAuth2DynamicClientRegistrationRequest {
	m := api.OAuth2DynamicClientRegistrationRequest{
		Provider:                  int32(d.Get("oauth2_provider").(int)),
		DefaultApplicationGroup:   helpers.GetP[string](d, "default_application_group"),
		OverrideAuthorizationFlow: *api.NewNullableString(helpers.GetP[string](d, "override_authorization_flow")),
		OverrideInvalidationFlow:  *api.NewNullableString(helpers.GetP[string](d, "override_invalidation_flow")),
		OverridePropertyMappings:  helpers.CastSlice[string](d, "override_property_mappings"),
		AccessTokenValidity:       helpers.GetP[string](d, "access_token_validity"),
		RefreshTokenValidity:      helpers.GetP[string](d, "refresh_token_validity"),
		AllowedGrantTypes:         helpers.CastSliceString[api.GrantTypeEnum](helpers.CastSlice[string](d, "allowed_grant_types")),
	}
	if pem, ok := d.GetOk("policy_engine_mode"); ok {
		m.PolicyEngineMode = api.PolicyEngineMode(pem.(string)).Ptr()
	}
	return &m
}

func resourceProviderOAuth2DCRCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceProviderOAuth2DCRSchemaToModel(d)

	res, hr, err := c.client.ProvidersAPI.ProvidersOauth2DcrCreate(ctx).OAuth2DynamicClientRegistrationRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceProviderOAuth2DCRRead(ctx, d, m)
}

func resourceProviderOAuth2DCRRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.ProvidersAPI.ProvidersOauth2DcrRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "oauth2_provider", res.Provider)
	helpers.SetWrapper(d, "default_application_group", res.DefaultApplicationGroup)
	helpers.SetWrapper(d, "override_authorization_flow", res.OverrideAuthorizationFlow.Get())
	helpers.SetWrapper(d, "override_invalidation_flow", res.OverrideInvalidationFlow.Get())
	helpers.SetWrapper(d, "override_property_mappings", res.OverridePropertyMappings)
	helpers.SetWrapper(d, "access_token_validity", res.AccessTokenValidity)
	helpers.SetWrapper(d, "refresh_token_validity", res.RefreshTokenValidity)
	helpers.SetWrapper(d, "allowed_grant_types", res.AllowedGrantTypes)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	return diags
}

func resourceProviderOAuth2DCRUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceProviderOAuth2DCRSchemaToModel(d)

	res, hr, err := c.client.ProvidersAPI.ProvidersOauth2DcrUpdate(ctx, d.Id()).OAuth2DynamicClientRegistrationRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceProviderOAuth2DCRRead(ctx, d, m)
}

func resourceProviderOAuth2DCRDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.ProvidersAPI.ProvidersOauth2DcrDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
