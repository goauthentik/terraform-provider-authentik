package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScopeMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScopeMappingCreate,
		ReadContext:   resourceScopeMappingRead,
		UpdateContext: resourceScopeMappingUpdate,
		DeleteContext: resourceScopeMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scope_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceScopeMappingSchemaToProvider(d *schema.ResourceData) (*api.ScopeMappingRequest, diag.Diagnostics) {
	r := api.ScopeMappingRequest{
		Name:       d.Get("name").(string),
		ScopeName:  d.Get("scope_name").(string),
		Expression: d.Get("expression").(string),
	}
	if de, dSet := d.GetOk("description"); dSet {
		r.Description = stringToPointer(de.(string))
	}
	return &r, nil
}

func resourceScopeMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceScopeMappingSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeCreate(ctx).ScopeMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceScopeMappingRead(ctx, d, m)
}

func resourceScopeMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("expression", res.Expression)
	d.Set("scope_name", res.ScopeName)
	d.Set("description", res.Description)
	return diags
}

func resourceScopeMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceScopeMappingSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeUpdate(ctx, d.Id()).ScopeMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceScopeMappingRead(ctx, d, m)
}

func resourceScopeMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsScopeDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
