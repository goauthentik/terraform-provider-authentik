package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceNotificationPropertyMapping() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourceNotificationPropertyMappingCreate,
		ReadContext:   resourceNotificationPropertyMappingRead,
		UpdateContext: resourceNotificationPropertyMappingUpdate,
		DeleteContext: resourceNotificationPropertyMappingDelete,
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

func resourceNotificationPropertyMappingSchemaToProvider(d *schema.ResourceData) *api.NotificationWebhookMappingRequest {
	r := api.NotificationWebhookMappingRequest{
		Name:       d.Get("name").(string),
		Expression: d.Get("expression").(string),
	}
	return &r
}

func resourceNotificationPropertyMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceNotificationPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationCreate(ctx).NotificationWebhookMappingRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceNotificationPropertyMappingRead(ctx, d, m)
}

func resourceNotificationPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceNotificationPropertyMappingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceNotificationPropertyMappingSchemaToProvider(d)

	res, hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationUpdate(ctx, d.Id()).NotificationWebhookMappingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceNotificationPropertyMappingRead(ctx, d, m)
}

func resourceNotificationPropertyMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PropertymappingsApi.PropertymappingsNotificationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
