package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingSourceSCIM() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage SCIM Source Property mappings",
		CreateContext: resourcePropertyMappingSourceSCIMCreate,
		ReadContext:   resourcePropertyMappingSourceSCIMRead,
		UpdateContext: resourcePropertyMappingSourceSCIMUpdate,
		DeleteContext: resourcePropertyMappingSourceSCIMDelete,
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
				DiffSuppressFunc: helpers.DiffSuppressExpression,
			},
		},
	}
}

func resourcePropertyMappingSourceSCIMSchemaToProvider(d *schema.ResourceData) *api.SCIMSourcePropertyMappingRequest {
	r := api.SCIMSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceSCIMCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceSCIMSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceScimCreate(ctx).SCIMSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSCIMRead(ctx, d, m)
}

func resourcePropertyMappingSourceSCIMRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceScimRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceSCIMUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceSCIMSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceScimUpdate(ctx, d.Id()).SCIMSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSCIMRead(ctx, d, m)
}

func resourcePropertyMappingSourceSCIMDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceScimDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
