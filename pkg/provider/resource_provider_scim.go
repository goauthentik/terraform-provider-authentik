package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderSCIM() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderSCIMCreate,
		ReadContext:   resourceProviderSCIMRead,
		UpdateContext: resourceProviderSCIMUpdate,
		DeleteContext: resourceProviderSCIMDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"auth_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SCIMAUTHENTICATIONMODEENUM_TOKEN,
				Description:      helpers.EnumToDescription(api.AllowedSCIMAuthenticationModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSCIMAuthenticationModeEnumEnumValues),
			},
			"auth_oauth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Slug of an OAuth source used for authentication",
			},
			"auth_oauth_params": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"compatibility_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.COMPATIBILITYMODEENUM_DEFAULT,
				Description:      helpers.EnumToDescription(api.AllowedCompatibilityModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedCompatibilityModeEnumEnumValues),
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
			"exclude_users_service_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_provider_config_cache_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "hours=1",
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"sync_page_timeout": {
				Type:             schema.TypeString,
				Default:          "minutes=30",
				Optional:         true,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"sync_page_size": {
				Type:     schema.TypeInt,
				Default:  100,
				Optional: true,
			},
		},
	}
}

func resourceProviderSCIMSchemaToProvider(d *schema.ResourceData) (*api.SCIMProviderRequest, diag.Diagnostics) {
	r := api.SCIMProviderRequest{
		Name:                              d.Get("name").(string),
		Url:                               d.Get("url").(string),
		AuthMode:                          helpers.CastString[api.SCIMAuthenticationModeEnum](helpers.GetP[string](d, "auth_mode")),
		AuthOauth:                         *api.NewNullableString(helpers.GetP[string](d, "auth_oauth")),
		Token:                             helpers.GetP[string](d, "token"),
		PropertyMappings:                  helpers.CastSlice[string](d, "property_mappings"),
		PropertyMappingsGroup:             helpers.CastSlice[string](d, "property_mappings_group"),
		ExcludeUsersServiceAccount:        new(d.Get("exclude_users_service_account").(bool)),
		CompatibilityMode:                 api.CompatibilityModeEnum(d.Get("compatibility_mode").(string)).Ptr(),
		FilterGroup:                       *api.NewNullableString(helpers.GetP[string](d, "filter_group")),
		DryRun:                            new(d.Get("dry_run").(bool)),
		ServiceProviderConfigCacheTimeout: helpers.GetP[string](d, "service_provider_config_cache_timeout"),
		SyncPageTimeout:                   helpers.GetP[string](d, "sync_page_timeout"),
		SyncPageSize:                      helpers.GetIntP(d, "sync_page_size"),
	}

	attr, err := helpers.GetJSON[map[string]any](d, "auth_oauth_params")
	r.AuthOauthParams = attr
	return &r, err
}

func resourceProviderSCIMCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderSCIMSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersScimCreate(ctx).SCIMProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSCIMRead(ctx, d, m)
}

func resourceProviderSCIMRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersScimRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "url", res.Url)
	helpers.SetWrapper(d, "token", res.Token)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.PropertyMappingsGroup,
	))
	helpers.SetWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	helpers.SetWrapper(d, "filter_group", res.FilterGroup.Get())
	helpers.SetWrapper(d, "dry_run", res.DryRun)
	helpers.SetWrapper(d, "compatibility_mode", res.CompatibilityMode)
	helpers.SetWrapper(d, "auth_mode", res.AuthMode)
	helpers.SetWrapper(d, "auth_oauth", res.AuthOauth.Get())
	helpers.SetWrapper(d, "service_provider_config_cache_timeout", res.ServiceProviderConfigCacheTimeout)
	helpers.SetWrapper(d, "sync_page_timeout", res.SyncPageTimeout)
	helpers.SetWrapper(d, "sync_page_size", res.SyncPageSize)
	return helpers.SetJSON(d, "auth_oauth_params", res.AuthOauthParams)
}

func resourceProviderSCIMUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app, diags := resourceProviderSCIMSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersScimUpdate(ctx, int32(id)).SCIMProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSCIMRead(ctx, d, m)
}

func resourceProviderSCIMDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersScimDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
