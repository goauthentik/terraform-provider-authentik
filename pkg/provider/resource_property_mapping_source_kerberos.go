package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingSourceKerberos() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Kerberos Source Property mappings",
		CreateContext: resourcePropertyMappingSourceKerberosCreate,
		ReadContext:   resourcePropertyMappingSourceKerberosRead,
		UpdateContext: resourcePropertyMappingSourceKerberosUpdate,
		DeleteContext: resourcePropertyMappingSourceKerberosDelete,
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

func resourcePropertyMappingSourceKerberosSchemaToProvider(d *schema.ResourceData) *api.KerberosSourcePropertyMappingRequest {
	r := api.KerberosSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceKerberosCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceKerberosSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceKerberosCreate(ctx).KerberosSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceKerberosRead(ctx, d, m)
}

func resourcePropertyMappingSourceKerberosRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceKerberosRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceKerberosUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceKerberosSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceKerberosUpdate(ctx, d.Id()).KerberosSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceKerberosRead(ctx, d, m)
}

func resourcePropertyMappingSourceKerberosDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceKerberosDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
