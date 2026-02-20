package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingProviderSAML() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage SAML Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderSAMLCreate,
		ReadContext:   resourcePropertyMappingProviderSAMLRead,
		UpdateContext: resourcePropertyMappingProviderSAMLUpdate,
		DeleteContext: resourcePropertyMappingProviderSAMLDelete,
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
				DiffSuppressFunc: helpers.DiffSuppressExpression,
			},
		},
	}
}

func resourcePropertyMappingProviderSAMLSchemaToProvider(d *schema.ResourceData) *api.SAMLPropertyMappingRequest {
	r := api.SAMLPropertyMappingRequest{
		Name:         d.Get("name").(string),
		SamlName:     d.Get("saml_name").(string),
		Expression:   d.Get("expression").(string),
		FriendlyName: *api.NewNullableString(helpers.GetP[string](d, "friendly_name")),
	}
	return &r
}

func resourcePropertyMappingProviderSAMLCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlCreate(ctx).SAMLPropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSAMLRead(ctx, d, m)
}

func resourcePropertyMappingProviderSAMLRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	helpers.SetWrapper(d, "saml_name", res.SamlName)
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName.Get())
	return diags
}

func resourcePropertyMappingProviderSAMLUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderSAMLSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlUpdate(ctx, d.Id()).SAMLPropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderSAMLRead(ctx, d, m)
}

func resourcePropertyMappingProviderSAMLDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderSamlDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
