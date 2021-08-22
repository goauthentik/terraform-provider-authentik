package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSAMLPropertyMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSAMLPropertyMappingCreate,
		ReadContext:   resourceSAMLPropertyMappingRead,
		UpdateContext: resourceSAMLPropertyMappingUpdate,
		DeleteContext: resourceSAMLPropertyMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"saml_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expression": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
		},
	}
}

func resourceSAMLPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.SAMLPropertyMappingRequest {
	r := api.SAMLPropertyMappingRequest{
		Name:       d.Get("name").(string),
		SamlName:   d.Get("saml_name").(string),
		Expression: d.Get("expression").(string),
	}
	if de, dSet := d.GetOk("friendly_name"); dSet {
		r.FriendlyName.Set(stringToPointer(de.(string)))
	}
	return &r
}

func resourceSAMLPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSAMLPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSamlCreate(ctx).SAMLPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceSAMLPropertyMappingRead(ctx, d, m)
}

func resourceSAMLPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("expression", res.Expression)
	d.Set("saml_name", res.SamlName)
	if res.FriendlyName.IsSet() {
		d.Set("friendly_name", res.FriendlyName.Get())
	}
	return diags
}

func resourceSAMLPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceSAMLPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSamlUpdate(ctx, d.Id()).SAMLPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceSAMLPropertyMappingRead(ctx, d, m)
}

func resourceSAMLPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
