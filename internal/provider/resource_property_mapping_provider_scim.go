package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingProviderSCIM() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage SCIM Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderSCIMCreate,
		ReadContext:   resourcePropertyMappingProviderSCIMRead,
		UpdateContext: resourcePropertyMappingProviderSCIMUpdate,
		DeleteContext: resourcePropertyMappingProviderSCIMDelete,
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
				Required:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
		},
	}
}

func resourcePropertyMappingProviderSCIMSchemaToProvider(d *schema.ResourceData) *api.SCIMMappingRequest {
	r := api.SCIMMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingProviderSCIMCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderSCIMSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderScimCreate(ctx).SCIMMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSCIMRead(ctx, d, m)
}

func resourcePropertyMappingProviderSCIMRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderScimRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingProviderSCIMUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderSCIMSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderScimUpdate(ctx, d.Id()).SCIMMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSCIMRead(ctx, d, m)
}

func resourcePropertyMappingProviderSCIMDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderScimDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
