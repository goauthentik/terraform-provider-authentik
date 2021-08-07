package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/imdario/mergo"
)

func resourceOutpost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOutpostCreate,
		ReadContext:   resourceOutpostRead,
		UpdateContext: resourceOutpostUpdate,
		DeleteContext: resourceOutpostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTPOSTTYPEENUM_PROXY,
			},
			"protocol_providers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"service_connection": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceOutpostSchemaToModel(d *schema.ResourceData, c *APIClient) (*api.OutpostRequest, diag.Diagnostics) {
	m := api.OutpostRequest{
		Name: d.Get("name").(string),
	}

	protocolProviders := d.Get("protocol_providers").([]interface{})
	m.Providers = make([]int32, len(protocolProviders))
	for i, prov := range protocolProviders {
		m.Providers[i] = int32(prov.(int))
	}

	if l, ok := d.Get("service_connection").(string); ok {
		m.ServiceConnection.Set(&l)
	} else {
		m.ServiceConnection.Set(nil)
	}

	defaultConfig, hr, err := c.client.OutpostsApi.OutpostsInstancesDefaultSettingsRetrieve(context.Background()).Execute()
	if err != nil {
		return nil, httpToDiag(hr, err)
	}
	m.Config = defaultConfig.Config
	if l, ok := d.Get("config").(string); ok {
		if l != "" {
			var c map[string]interface{}
			err := json.NewDecoder(strings.NewReader(l)).Decode(&c)
			if err != nil {
				return nil, diag.FromErr(err)
			}
			err = mergo.Merge(&m.Config, c)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}
	}

	m.Type = api.OutpostTypeEnum(d.Get("type").(string))
	return &m, nil
}

func resourceOutpostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceOutpostSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesCreate(ctx).OutpostRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("type", res.Type)
	d.Set("protocol_providers", res.Providers)
	if res.ServiceConnection.IsSet() {
		d.Set("service_connection", res.ServiceConnection.Get())
	}
	b, err := json.Marshal(res.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("config", string(b))
	return diags
}

func resourceOutpostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceOutpostSchemaToModel(d, c)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesUpdate(ctx, d.Id()).OutpostRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.OutpostsApi.OutpostsInstancesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
