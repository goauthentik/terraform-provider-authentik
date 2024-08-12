package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingSourcePlex() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Plex Source Property mappings",
		CreateContext: resourcePropertyMappingSourcePlexCreate,
		ReadContext:   resourcePropertyMappingSourcePlexRead,
		UpdateContext: resourcePropertyMappingSourcePlexUpdate,
		DeleteContext: resourcePropertyMappingSourcePlexDelete,
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

func resourcePropertyMappingSourcePlexSchemaToProvider(d *schema.ResourceData) *api.PlexSourcePropertyMappingRequest {
	r := api.PlexSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourcePlexCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourcePlexSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourcePlexCreate(ctx).PlexSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourcePlexRead(ctx, d, m)
}

func resourcePropertyMappingSourcePlexRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourcePlexRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourcePlexUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourcePlexSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourcePlexUpdate(ctx, d.Id()).PlexSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourcePlexRead(ctx, d, m)
}

func resourcePropertyMappingSourcePlexDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourcePlexDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
