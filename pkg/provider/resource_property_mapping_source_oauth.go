package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingSourceOAuth() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage OAuth Source Property mappings",
		CreateContext: resourcePropertyMappingSourceOAuthCreate,
		ReadContext:   resourcePropertyMappingSourceOAuthRead,
		UpdateContext: resourcePropertyMappingSourceOAuthUpdate,
		DeleteContext: resourcePropertyMappingSourceOAuthDelete,
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

func resourcePropertyMappingSourceOAuthSchemaToProvider(d *schema.ResourceData) *api.OAuthSourcePropertyMappingRequest {
	r := api.OAuthSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceOAuthCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceOAuthSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceOauthCreate(ctx).OAuthSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceOAuthRead(ctx, d, m)
}

func resourcePropertyMappingSourceOAuthRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceOauthRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceOAuthUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceOAuthSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsSourceOauthUpdate(ctx, d.Id()).OAuthSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceOAuthRead(ctx, d, m)
}

func resourcePropertyMappingSourceOAuthDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsSourceOauthDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
