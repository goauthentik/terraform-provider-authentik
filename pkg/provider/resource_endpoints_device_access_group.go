package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceEndpointsDeviceAccessGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Endpoint Devices --- ",
		CreateContext: resourceEndpointsDeviceAccessGroupCreate,
		ReadContext:   resourceEndpointsDeviceAccessGroupRead,
		UpdateContext: resourceEndpointsDeviceAccessGroupUpdate,
		DeleteContext: resourceEndpointsDeviceAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceEndpointsDeviceAccessGroupSchemaToModel(d *schema.ResourceData) (*api.DeviceAccessGroupRequest, diag.Diagnostics) {
	m := api.DeviceAccessGroupRequest{
		Name: d.Get("name").(string),
	}
	return &m, nil
}

func resourceEndpointsDeviceAccessGroupCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEndpointsDeviceAccessGroupSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsApi.EndpointsDeviceAccessGroupsCreate(ctx).DeviceAccessGroupRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceEndpointsDeviceAccessGroupRead(ctx, d, m)
}

func resourceEndpointsDeviceAccessGroupRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EndpointsApi.EndpointsDeviceAccessGroupsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	return diags
}

func resourceEndpointsDeviceAccessGroupUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceEndpointsDeviceAccessGroupSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.EndpointsApi.EndpointsDeviceAccessGroupsUpdate(ctx, d.Id()).DeviceAccessGroupRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.PbmUuid)
	return resourceEndpointsDeviceAccessGroupRead(ctx, d, m)
}

func resourceEndpointsDeviceAccessGroupDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EndpointsApi.EndpointsDeviceAccessGroupsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
