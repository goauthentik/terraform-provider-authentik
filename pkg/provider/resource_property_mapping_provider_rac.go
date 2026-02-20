package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingProviderRAC() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage RAC Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderRACCreate,
		ReadContext:   resourcePropertyMappingProviderRACRead,
		UpdateContext: resourcePropertyMappingProviderRACUpdate,
		DeleteContext: resourcePropertyMappingProviderRACDelete,
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
				Optional:         true,
				DiffSuppressFunc: helpers.DiffSuppressExpression,
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourcePropertyMappingProviderRACSchemaToProvider(d *schema.ResourceData) (*api.RACPropertyMappingRequest, diag.Diagnostics) {
	r := api.RACPropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: helpers.GetP[string](d, "expression"),
	}

	settings, err := helpers.GetJSON[map[string]any](d, ("settings"))
	r.StaticSettings = settings
	return &r, err
}

func resourcePropertyMappingProviderRACCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourcePropertyMappingProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacCreate(ctx).RACPropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRACRead(ctx, d, m)
}

func resourcePropertyMappingProviderRACRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.GetExpression())
	return helpers.SetJSON(d, "settings", res.StaticSettings)
}

func resourcePropertyMappingProviderRACUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourcePropertyMappingProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacUpdate(ctx, d.Id()).RACPropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRACRead(ctx, d, m)
}

func resourcePropertyMappingProviderRACDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
