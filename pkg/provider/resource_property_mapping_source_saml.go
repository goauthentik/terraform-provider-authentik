package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				DiffSuppressFunc: helpers.DiffSuppressExpression,
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

func resourcePropertyMappingSourceSAMLCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlCreate(ctx).SAMLSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSAMLRead(ctx, d, m)
}

func resourcePropertyMappingSourceSAMLRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceSAMLUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlUpdate(ctx, d.Id()).SAMLSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceSAMLRead(ctx, d, m)
}

func resourcePropertyMappingSourceSAMLDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
