package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
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
		Token:                      d.Get("token").(string),
		PropertyMappings:           castSlice[string](d.Get("property_mappings").([]interface{})),
		PropertyMappingsGroup:      castSlice[string](d.Get("property_mappings_group").([]interface{})),
		ExcludeUsersServiceAccount: api.PtrBool(d.Get("exclude_users_service_account").(bool)),
	}
	if l, ok := d.Get("filter_group").(string); ok {
		r.FilterGroup = *api.NewNullableString(&l)
	}
	return &r
}

func resourceProviderSCIMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderSCIMSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersScimCreate(ctx).SCIMProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "url", res.Url)
	setWrapper(d, "token", res.Token)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.PropertyMappings))
	localGroupMappings := castSlice[string](d.Get("property_mappings_group").([]interface{}))
	setWrapper(d, "property_mappings_group", listConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	setWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	setWrapper(d, "filter_group", res.FilterGroup.Get())
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
		return httpToDiag(d, hr, err)
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
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
