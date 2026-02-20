package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingProviderMicrosoftEntra() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Microsoft Entra Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderMicrosoftEntraCreate,
		ReadContext:   resourcePropertyMappingProviderMicrosoftEntraRead,
		UpdateContext: resourcePropertyMappingProviderMicrosoftEntraUpdate,
		DeleteContext: resourcePropertyMappingProviderMicrosoftEntraDelete,
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

func resourcePropertyMappingProviderMicrosoftEntraSchemaToProvider(d *schema.ResourceData) *api.MicrosoftEntraProviderMappingRequest {
	r := api.MicrosoftEntraProviderMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingProviderMicrosoftEntraCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderMicrosoftEntraSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraCreate(ctx).MicrosoftEntraProviderMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderMicrosoftEntraRead(ctx, d, m)
}

func resourcePropertyMappingProviderMicrosoftEntraRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingProviderMicrosoftEntraUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderMicrosoftEntraSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraUpdate(ctx, d.Id()).MicrosoftEntraProviderMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderMicrosoftEntraRead(ctx, d, m)
}

func resourcePropertyMappingProviderMicrosoftEntraDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
