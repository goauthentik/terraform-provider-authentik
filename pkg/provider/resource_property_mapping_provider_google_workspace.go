package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingProviderGoogleWorkspace() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Google Workspace Provider Property mappings",
		CreateContext: resourcePropertyMappingProviderGoogleWorkspaceCreate,
		ReadContext:   resourcePropertyMappingProviderGoogleWorkspaceRead,
		UpdateContext: resourcePropertyMappingProviderGoogleWorkspaceUpdate,
		DeleteContext: resourcePropertyMappingProviderGoogleWorkspaceDelete,
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

func resourcePropertyMappingProviderGoogleWorkspaceSchemaToProvider(d *schema.ResourceData) *api.GoogleWorkspaceProviderMappingRequest {
	r := api.GoogleWorkspaceProviderMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingProviderGoogleWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingProviderGoogleWorkspaceSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderGoogleWorkspaceCreate(ctx).GoogleWorkspaceProviderMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderGoogleWorkspaceRead(ctx, d, m)
}

func resourcePropertyMappingProviderGoogleWorkspaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderGoogleWorkspaceRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingProviderGoogleWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingProviderGoogleWorkspaceSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsProviderGoogleWorkspaceUpdate(ctx, d.Id()).GoogleWorkspaceProviderMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingProviderGoogleWorkspaceRead(ctx, d, m)
}

func resourcePropertyMappingProviderGoogleWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsProviderGoogleWorkspaceDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
