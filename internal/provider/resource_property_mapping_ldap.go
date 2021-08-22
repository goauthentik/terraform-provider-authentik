package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLDAPPropertyMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLDAPPropertyMappingCreate,
		ReadContext:   resourceLDAPPropertyMappingRead,
		UpdateContext: resourceLDAPPropertyMappingUpdate,
		DeleteContext: resourceLDAPPropertyMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"object_field": {
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

func resourceLDAPPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.LDAPPropertyMappingRequest {
	r := api.LDAPPropertyMappingRequest{
		Name:        d.Get("name").(string),
		ObjectField: d.Get("object_field").(string),
		Expression:  d.Get("expression").(string),
	}
	return &r
}

func resourceLDAPPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceLDAPPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsLdapCreate(ctx).LDAPPropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceLDAPPropertyMappingRead(ctx, d, m)
}

func resourceLDAPPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("expression", res.Expression)
	d.Set("object_field", res.ObjectField)
	return diags
}

func resourceLDAPPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceLDAPPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsLdapUpdate(ctx, d.Id()).LDAPPropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceLDAPPropertyMappingRead(ctx, d, m)
}

func resourceLDAPPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
