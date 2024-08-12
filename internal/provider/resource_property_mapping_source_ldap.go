package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingSourceLDAP() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage LDAP Source Property mappings",
		CreateContext: resourcePropertyMappingSourceLDAPCreate,
		ReadContext:   resourcePropertyMappingSourceLDAPRead,
		UpdateContext: resourcePropertyMappingSourceLDAPUpdate,
		DeleteContext: resourcePropertyMappingSourceLDAPDelete,
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

func resourcePropertyMappingSourceLDAPSchemaToProvider(d *schema.ResourceData) *api.LDAPSourcePropertyMappingRequest {
	r := api.LDAPSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceLDAPSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceLdapCreate(ctx).LDAPSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceLDAPRead(ctx, d, m)
}

func resourcePropertyMappingSourceLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceLDAPSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceLdapUpdate(ctx, d.Id()).LDAPSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceLDAPRead(ctx, d, m)
}

func resourcePropertyMappingSourceLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
