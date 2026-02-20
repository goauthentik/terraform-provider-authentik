package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceApplicationEntitlement() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceApplicationEntitlementCreate,
		ReadContext:   resourceApplicationEntitlementRead,
		UpdateContext: resourceApplicationEntitlementUpdate,
		DeleteContext: resourceApplicationEntitlementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourceApplicationEntitlementSchemaToModel(d *schema.ResourceData) (*api.ApplicationEntitlementRequest, diag.Diagnostics) {
	m := api.ApplicationEntitlementRequest{
		Name: d.Get("name").(string),
		App:  d.Get("application").(string),
	}

	attr, err := helpers.GetJSON[map[string]any](d, ("attributes"))
	m.Attributes = attr
	return &m, err
}

func resourceApplicationEntitlementCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceApplicationEntitlementSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreApplicationEntitlementsCreate(ctx).ApplicationEntitlementRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceApplicationEntitlementRead(ctx, d, m)
}

func resourceApplicationEntitlementRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreApplicationEntitlementsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "application", res.App)
	return helpers.SetJSON(d, "attributes", res.Attributes)
}

func resourceApplicationEntitlementUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceApplicationEntitlementSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreApplicationEntitlementsUpdate(ctx, d.Id()).ApplicationEntitlementRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceApplicationEntitlementRead(ctx, d, m)
}

func resourceApplicationEntitlementDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreApplicationEntitlementsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
