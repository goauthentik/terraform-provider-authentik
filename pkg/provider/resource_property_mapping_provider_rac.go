package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				DiffSuppressFunc: diffSuppressExpression,
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      JSONDescription,
				DiffSuppressFunc: diffSuppressJSON,
				ValidateDiagFunc: ValidateJSON,
			},
		},
	}
}

func resourcePropertyMappingProviderRACSchemaToProvider(d *schema.ResourceData) (*api.RACPropertyMappingRequest, diag.Diagnostics) {
	r := api.RACPropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: getP[string](d, "expression"),
	}

	settings, err := getJSON[map[string]interface{}](d, ("settings"))
	r.StaticSettings = settings
	return &r, err
}

func resourcePropertyMappingProviderRACCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourcePropertyMappingProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacCreate(ctx).RACPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRACRead(ctx, d, m)
}

func resourcePropertyMappingProviderRACRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.GetExpression())
	b, err := json.Marshal(res.StaticSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "settings", string(b))
	return diags
}

func resourcePropertyMappingProviderRACUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourcePropertyMappingProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacUpdate(ctx, d.Id()).RACPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderRACRead(ctx, d, m)
}

func resourcePropertyMappingProviderRACDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderRacDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
