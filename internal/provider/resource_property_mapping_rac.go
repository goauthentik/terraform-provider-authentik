package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceRACPropertyMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceRACPropertyMappingCreate,
		ReadContext:   resourceRACPropertyMappingRead,
		UpdateContext: resourceRACPropertyMappingUpdate,
		DeleteContext: resourceRACPropertyMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      "JSON format expected. Use jsonencode() to pass objects.",
				DiffSuppressFunc: diffSuppressJSON,
			},
		},
	}
}

func resourceRACPropertyMappingSchemaToProvider(d *schema.ResourceData) (*api.RACPropertyMappingRequest, diag.Diagnostics) {
	r := api.RACPropertyMappingRequest{
		Name: d.Get("name").(string),
	}
	if s, sok := d.GetOk("expression"); sok && s.(string) != "" {
		r.Expression = api.PtrString(s.(string))
	}

	attr := make(map[string]interface{})
	if l, ok := d.Get("settings").(string); ok && l != "" {
		err := json.NewDecoder(strings.NewReader(l)).Decode(&attr)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}
	r.StaticSettings = attr
	return &r, nil
}

func resourceRACPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceRACPropertyMappingSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRacCreate(ctx).RACPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRACPropertyMappingRead(ctx, d, m)
}

func resourceRACPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRacRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.GetExpression())
	b, err := json.Marshal(res.StaticSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "settings", string(b))
	return diags
}

func resourceRACPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceRACPropertyMappingSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRacUpdate(ctx, d.Id()).RACPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRACPropertyMappingRead(ctx, d, m)
}

func resourceRACPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsRacDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
