package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceEndpointsEnrollmentToken() *schema.Resource {
	return &schema.Resource{
		Description:   "Endpoint Devices --- ",
		CreateContext: resourceEndpointsEnrollmentTokenCreate,
		ReadContext:   resourceEndpointsEnrollmentTokenRead,
		UpdateContext: resourceEndpointsEnrollmentTokenUpdate,
		DeleteContext: resourceEndpointsEnrollmentTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Computed
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"expires_in": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Meta
			"retrieve_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Actual
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device_access_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connector": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiring": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceEndpointsEnrollmentTokenSchemaToModel(d *schema.ResourceData) (*api.EnrollmentTokenRequest, diag.Diagnostics) {
	m := api.EnrollmentTokenRequest{
		Name:        d.Get("name").(string),
		Expiring:    new(d.Get("expiring").(bool)),
		DeviceGroup: *api.NewNullableString(helpers.GetP[string](d, "device_access_group")),
		Connector:   d.Get("connector").(string),
	}

	if l, ok := d.Get("expires").(string); ok && l != "" {
		t, err := time.Parse(time.RFC3339, l)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		m.Expires.Set(&t)
	}
	return &m, nil
}

func resourceEndpointsEnrollmentTokenCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEndpointsEnrollmentTokenSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsApi.EndpointsAgentsEnrollmentTokensCreate(ctx).EnrollmentTokenRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.TokenUuid)
	return resourceEndpointsEnrollmentTokenRead(ctx, d, m)
}

func resourceEndpointsEnrollmentTokenRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EndpointsApi.EndpointsAgentsEnrollmentTokensRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "connector", res.Connector)
	helpers.SetWrapper(d, "device_access_group", res.DeviceGroup.Get())
	if res.Expires.IsSet() && res.Expires.Get() != nil {
		helpers.SetWrapper(d, "expires_in", time.Until(*res.Expires.Get()).Seconds())
	}
	if rt, ok := d.Get("retrieve_key").(bool); ok && rt {
		res, hr, err := c.client.EndpointsApi.EndpointsAgentsEnrollmentTokensViewKeyRetrieve(ctx, d.Id()).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		helpers.SetWrapper(d, "key", res.Key)
	}
	return diags
}

func resourceEndpointsEnrollmentTokenUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceEndpointsEnrollmentTokenSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.EndpointsApi.EndpointsAgentsEnrollmentTokensUpdate(ctx, d.Id()).EnrollmentTokenRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.TokenUuid)
	return resourceEndpointsEnrollmentTokenRead(ctx, d, m)
}

func resourceEndpointsEnrollmentTokenDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EndpointsApi.EndpointsAgentsEnrollmentTokensDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
