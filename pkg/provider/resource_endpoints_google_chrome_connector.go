package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceEndpointsGoogleChromeConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "Endpoint Devices --- ",
		CreateContext: resourceEndpointsGoogleChromeConnectorCreate,
		ReadContext:   resourceEndpointsGoogleChromeConnectorRead,
		UpdateContext: resourceEndpointsGoogleChromeConnectorUpdate,
		DeleteContext: resourceEndpointsGoogleChromeConnectorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"credentials": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"chrome_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEndpointsGoogleChromeConnectorSchemaToProvider(d *schema.ResourceData) (*api.GoogleChromeConnectorRequest, diag.Diagnostics) {
	r := api.GoogleChromeConnectorRequest{
		Name:    d.Get("name").(string),
		Enabled: new(d.Get("enabled").(bool)),
	}

	credentials, err := helpers.GetJSON[map[string]any](d, "credentials")
	r.Credentials = credentials
	return &r, err
}

func resourceEndpointsGoogleChromeConnectorCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceEndpointsGoogleChromeConnectorSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsAPI.EndpointsGoogleChromeConnectorsCreate(ctx).GoogleChromeConnectorRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(*res.ConnectorUuid)
	return resourceEndpointsGoogleChromeConnectorRead(ctx, d, m)
}

func resourceEndpointsGoogleChromeConnectorRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.EndpointsAPI.EndpointsGoogleChromeConnectorsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "chrome_url", res.ChromeUrl.Get())
	return helpers.SetJSON(d, "credentials", res.Credentials)
}

func resourceEndpointsGoogleChromeConnectorUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEndpointsGoogleChromeConnectorSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsAPI.EndpointsGoogleChromeConnectorsUpdate(ctx, d.Id()).GoogleChromeConnectorRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(*res.ConnectorUuid)
	return resourceEndpointsGoogleChromeConnectorRead(ctx, d, m)
}

func resourceEndpointsGoogleChromeConnectorDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EndpointsAPI.EndpointsGoogleChromeConnectorsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
