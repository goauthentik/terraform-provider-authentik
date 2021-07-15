package authentik

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_provider": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"meta_launch_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_icon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_publisher": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},
		},
	}
}

func resourceApplicationSchemaToModel(d *schema.ResourceData) (*api.ApplicationRequest, diag.Diagnostics) {
	m := api.ApplicationRequest{
		Name:     d.Get("name").(string),
		Slug:     d.Get("slug").(string),
		Provider: api.NullableInt32{},
	}

	if p, pSet := d.GetOk("protocol_provider"); pSet {
		i := int32(p.(int))
		m.Provider.Set(&i)
	} else {
		m.Provider.Set(nil)
	}

	if l, ok := d.Get("meta_launch_url").(string); ok {
		m.MetaLaunchUrl = &l
	}
	if l, ok := d.Get("meta_description").(string); ok {
		m.MetaDescription = &l
	}
	if l, ok := d.Get("meta_publisher").(string); ok {
		m.MetaPublisher = &l
	}

	pm := d.Get("policy_engine_mode").(string)
	var pma api.PolicyEngineMode
	switch pm {
	case string(api.POLICYENGINEMODE_ALL):
		pma = api.POLICYENGINEMODE_ALL
	case string(api.POLICYENGINEMODE_ANY):
		pma = api.POLICYENGINEMODE_ANY
	default:
		return nil, diag.Errorf("invalid policy_engine_mode %s", pm)
	}
	m.PolicyEngineMode = &pma
	return &m, nil
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	app, diags := resourceApplicationSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreApplicationsCreate(ctx).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Slug)
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)

	res, hr, err := c.client.CoreApi.CoreApplicationsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	if prov := res.Provider.Get(); prov != nil {
		d.Set("protocol_provider", int(*prov))
	}
	d.Set("meta_launch_url", res.MetaLaunchUrl)
	d.Set("meta_description", res.MetaDescription)
	d.Set("meta_publisher", res.MetaPublisher)
	d.Set("policy_engine_mode", res.PolicyEngineMode)
	return diags
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	app, di := resourceApplicationSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.CoreApi.CoreApplicationsUpdate(ctx, d.Id()).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Slug)
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	hr, err := c.client.CoreApi.CoreApplicationsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
