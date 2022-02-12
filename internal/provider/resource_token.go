package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTokenCreate,
		ReadContext:   resourceTokenRead,
		UpdateContext: resourceTokenUpdate,
		DeleteContext: resourceTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Computed
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Meta
			"retrieve_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Actual
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"intent": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.INTENTENUM_API,
			},
			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTokenSchemaToModel(d *schema.ResourceData, c *APIClient) (*api.TokenRequest, diag.Diagnostics) {
	m := api.TokenRequest{
		Identifier: d.Get("identifier").(string),
		User:       intToPointer(d.Get("user").(int)),
		Expiring:   boolToPointer(d.Get("expiring").(bool)),
	}

	if l, ok := d.Get("description").(string); ok {
		m.Description = &l
	}

	if l, ok := d.Get("expires").(string); ok && l != "" {
		t, err := time.Parse(time.RFC3339, l)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		m.Expires = &t
	}
	int := api.IntentEnum(d.Get("intent").(string))
	m.Intent = &int
	return &m, nil
}

func resourceTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceTokenSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreTokensCreate(ctx).TokenRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Identifier)
	return resourceTokenRead(ctx, d, m)
}

func resourceTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreTokensRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("identifier", res.Identifier)
	d.Set("user", res.User)
	d.Set("description", res.Description)
	d.Set("intent", res.Intent)
	d.Set("expires_in", time.Until(*res.Expires).Seconds())
	if rt, ok := d.Get("retrieve_key").(bool); ok && rt {
		res, hr, err := c.client.CoreApi.CoreTokensViewKeyRetrieve(ctx, d.Id()).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
		d.Set("key", res.Key)
	}
	return diags
}

func resourceTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceTokenSchemaToModel(d, c)
	if di != nil {
		return di
	}
	res, hr, err := c.client.CoreApi.CoreTokensUpdate(ctx, d.Id()).TokenRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Identifier)
	return resourceTokenRead(ctx, d, m)
}

func resourceTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreTokensDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
