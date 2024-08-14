package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePropertyMappingNotification() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- Manage Notification Property mappings",
		CreateContext: resourcePropertyMappingNotificationCreate,
		ReadContext:   resourcePropertyMappingNotificationRead,
		UpdateContext: resourcePropertyMappingNotificationUpdate,
		DeleteContext: resourcePropertyMappingNotificationDelete,
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

func resourcePropertyMappingNotificationSchemaToProvider(d *schema.ResourceData) *api.NotificationWebhookMappingRequest {
	r := api.NotificationWebhookMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourcePropertyMappingNotificationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePropertyMappingNotificationSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationCreate(ctx).NotificationWebhookMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingNotificationRead(ctx, d, m)
}

func resourcePropertyMappingNotificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePropertyMappingNotificationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePropertyMappingNotificationSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationUpdate(ctx, d.Id()).NotificationWebhookMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePropertyMappingNotificationRead(ctx, d, m)
}

func resourcePropertyMappingNotificationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
