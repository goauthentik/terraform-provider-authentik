package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingProviderRadius() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Radius Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderRadiusCreate,
		ReadContext:   resourcePropertyMappingProviderRadiusRead,
		UpdateContext: resourcePropertyMappingProviderRadiusUpdate,
		DeleteContext: resourcePropertyMappingProviderRadiusDelete,
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

func resourcePropertyMappingProviderRadiusSchemaToProvider(d *schema.ResourceData) *api.RadiusProviderPropertyMappingRequest {
	r := api.RadiusProviderPropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingProviderRadiusCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRadiusCreate(ctx).RadiusProviderPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRadiusRead(ctx, d, m)
}

func resourcePropertyMappingProviderRadiusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRadiusRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingProviderRadiusUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRadiusUpdate(ctx, d.Id()).RadiusProviderPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRadiusRead(ctx, d, m)
}

func resourcePropertyMappingProviderRadiusDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRadiusDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
