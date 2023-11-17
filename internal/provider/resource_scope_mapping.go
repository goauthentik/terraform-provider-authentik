package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceScopeMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
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
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
		},
	}
}

func resourceScopeMappingSchemaToProvider(d *schema.ResourceData) *api.ScopeMappingRequest {
	r := api.ScopeMappingRequest{
		Name:       d.Get("name").(string),
		ScopeName:  d.Get("scope_name").(string),
		Expression: d.Get("expression").(string),
	}
	if de, dSet := d.GetOk("description"); dSet {
		r.Description = api.PtrString(de.(string))
	}
	return &r
}

func resourceScopeMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceScopeMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeCreate(ctx).ScopeMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceScopeMappingRead(ctx, d, m)
}

func resourceScopeMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	setWrapper(d, "scope_name", res.ScopeName)
	setWrapper(d, "description", res.Description)
	return diags
}

func resourceScopeMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceScopeMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScopeUpdate(ctx, d.Id()).ScopeMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceScopeMappingRead(ctx, d, m)
}

func resourceScopeMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsScopeDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
