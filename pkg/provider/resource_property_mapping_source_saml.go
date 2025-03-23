package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingSourceSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage SAML Source Property mappings",
		CreateContext: resourcePropertyMappingSourceSAMLCreate,
		ReadContext:   resourcePropertyMappingSourceSAMLRead,
		UpdateContext: resourcePropertyMappingSourceSAMLUpdate,
		DeleteContext: resourcePropertyMappingSourceSAMLDelete,
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

func resourcePropertyMappingSourceSAMLSchemaToProvider(d *schema.ResourceData) *api.SAMLSourcePropertyMappingRequest {
	r := api.SAMLSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlCreate(ctx).SAMLSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSAMLRead(ctx, d, m)
}

func resourcePropertyMappingSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlUpdate(ctx, d.Id()).SAMLSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSAMLRead(ctx, d, m)
}

func resourcePropertyMappingSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
