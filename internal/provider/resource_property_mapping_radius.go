package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceRadiusProviderPropertyMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceRadiusProviderPropertyMappingCreate,
		ReadContext:   resourceRadiusProviderPropertyMappingRead,
		UpdateContext: resourceRadiusProviderPropertyMappingUpdate,
		DeleteContext: resourceRadiusProviderPropertyMappingDelete,
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

func resourceRadiusProviderPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.RadiusProviderPropertyMappingRequest {
	r := api.RadiusProviderPropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourceRadiusProviderPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceRadiusProviderPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRadiusCreate(ctx).RadiusProviderPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRadiusProviderPropertyMappingRead(ctx, d, m)
}

func resourceRadiusProviderPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRadiusRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourceRadiusProviderPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRadiusProviderPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsRadiusUpdate(ctx, d.Id()).RadiusProviderPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRadiusProviderPropertyMappingRead(ctx, d, m)
}

func resourceRadiusProviderPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsRadiusDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
