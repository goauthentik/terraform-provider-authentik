package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSCIMPropertyMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceSCIMPropertyMappingCreate,
		ReadContext:   resourceSCIMPropertyMappingRead,
		UpdateContext: resourceSCIMPropertyMappingUpdate,
		DeleteContext: resourceSCIMPropertyMappingDelete,
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

func resourceSCIMPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.SCIMMappingRequest {
	r := api.SCIMMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourceSCIMPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSCIMPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScimCreate(ctx).SCIMMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceSCIMPropertyMappingRead(ctx, d, m)
}

func resourceSCIMPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScimRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourceSCIMPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceSCIMPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsScimUpdate(ctx, d.Id()).SCIMMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceSCIMPropertyMappingRead(ctx, d, m)
}

func resourceSCIMPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsScimDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
