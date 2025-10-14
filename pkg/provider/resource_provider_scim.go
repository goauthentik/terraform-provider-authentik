package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
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
				Required:  true,
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
		},
	}
}

func resourceProviderSCIMSchemaToProvider(d *schema.ResourceData) *api.SCIMProviderRequest {
	r := api.SCIMProviderRequest{
		Name:                       d.Get("name").(string),
		Url:                        d.Get("url").(string),
		Token:                      helpers.GetP[string](d, "token"),
		PropertyMappings:           helpers.CastSlice_New[string](d, "property_mappings"),
		PropertyMappingsGroup:      helpers.CastSlice_New[string](d, "property_mappings_group"),
		ExcludeUsersServiceAccount: api.PtrBool(d.Get("exclude_users_service_account").(bool)),
		CompatibilityMode:          api.CompatibilityModeEnum(d.Get("compatibility_mode").(string)).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.GetP[string](d, "filter_group")),
		DryRun:                     api.PtrBool(d.Get("dry_run").(bool)),
	}
	return &r
}

func resourceProviderSCIMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderSCIMSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersScimCreate(ctx).SCIMProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSCIMRead(ctx, d, m)
}

func resourceProviderSCIMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
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
	localMappings := helpers.CastSlice[string](d.Get("property_mappings").([]interface{}))
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(localMappings, res.PropertyMappings))
	localGroupMappings := helpers.CastSlice[string](d.Get("property_mappings_group").([]interface{}))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	helpers.SetWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	helpers.SetWrapper(d, "filter_group", res.FilterGroup.Get())
	helpers.SetWrapper(d, "dry_run", res.DryRun)
	helpers.SetWrapper(d, "compatibility_mode", res.CompatibilityMode)
	return diags
}

func resourceProviderSCIMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderSCIMSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersScimUpdate(ctx, int32(id)).SCIMProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderSCIMRead(ctx, d, m)
}

func resourceProviderSCIMDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
