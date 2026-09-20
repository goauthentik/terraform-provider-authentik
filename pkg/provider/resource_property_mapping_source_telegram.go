package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePropertyMappingSourceTelegram() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Telegram Source Property mappings",
		CreateContext: resourcePropertyMappingSourceTelegramCreate,
		ReadContext:   resourcePropertyMappingSourceTelegramRead,
		UpdateContext: resourcePropertyMappingSourceTelegramUpdate,
		DeleteContext: resourcePropertyMappingSourceTelegramDelete,
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

func resourcePropertyMappingSourceTelegramSchemaToProvider(d *schema.ResourceData) *api.TelegramSourcePropertyMappingRequest {
	r := api.TelegramSourcePropertyMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingSourceTelegramCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingSourceTelegramSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsAPI.PropertymappingsSourceTelegramCreate(ctx).TelegramSourcePropertyMappingRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceTelegramRead(ctx, d, m)
}

func resourcePropertyMappingSourceTelegramRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsAPI.PropertymappingsSourceTelegramRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingSourceTelegramUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingSourceTelegramSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsAPI.PropertymappingsSourceTelegramUpdate(ctx, d.Id()).TelegramSourcePropertyMappingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingSourceTelegramRead(ctx, d, m)
}

func resourcePropertyMappingSourceTelegramDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsAPI.PropertymappingsSourceTelegramDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
