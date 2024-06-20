package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceMicrosoftEntraPropertyMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceMicrosoftEntraPropertyMappingCreate,
		ReadContext:   resourceMicrosoftEntraPropertyMappingRead,
		UpdateContext: resourceMicrosoftEntraPropertyMappingUpdate,
		DeleteContext: resourceMicrosoftEntraPropertyMappingDelete,
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

func resourceMicrosoftEntraPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.MicrosoftEntraProviderMappingRequest {
	r := api.MicrosoftEntraProviderMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourceMicrosoftEntraPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceMicrosoftEntraPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraCreate(ctx).MicrosoftEntraProviderMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceMicrosoftEntraPropertyMappingRead(ctx, d, m)
}

func resourceMicrosoftEntraPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourceMicrosoftEntraPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceMicrosoftEntraPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraUpdate(ctx, d.Id()).MicrosoftEntraProviderMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceMicrosoftEntraPropertyMappingRead(ctx, d, m)
}

func resourceMicrosoftEntraPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderMicrosoftEntraDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
