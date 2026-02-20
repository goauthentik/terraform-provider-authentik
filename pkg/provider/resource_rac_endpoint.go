package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRACEndpoint() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceRACEndpointCreate,
		ReadContext:   resourceRACEndpointRead,
		UpdateContext: resourceRACEndpointUpdate,
		DeleteContext: resourceRACEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_provider": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"protocol": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedProtocolEnumEnumValues),
				Description:      helpers.EnumToDescription(api.AllowedProtocolEnumEnumValues),
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"maximum_connections": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"settings": {
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

func resourceRACEndpointSchemaToProvider(d *schema.ResourceData) (*api.EndpointRequest, diag.Diagnostics) {
	r := api.EndpointRequest{
		Name:               d.Get("name").(string),
		Provider:           int32(d.Get("protocol_provider").(int)),
		Protocol:           api.ProtocolEnum(d.Get("protocol").(string)),
		Host:               d.Get("host").(string),
		AuthMode:           api.ENDPOINTAUTHMODEENUM_PROMPT,
		PropertyMappings:   helpers.CastSlice[string](d, "property_mappings"),
		MaximumConnections: new(int32(d.Get("maximum_connections").(int))),
	}

	attr, err := helpers.GetJSON[map[string]any](d, ("settings"))
	r.Settings = attr
	return &r, err
}

func resourceRACEndpointCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceRACEndpointSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.RacApi.RacEndpointsCreate(ctx).EndpointRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRACEndpointRead(ctx, d, m)
}

func resourceRACEndpointRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	res, hr, err := c.client.RacApi.RacEndpointsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "protocol_provider", res.Provider)
	helpers.SetWrapper(d, "host", res.Host)
	helpers.SetWrapper(d, "protocol", res.Protocol)
	helpers.SetWrapper(d, "maximum_connections", res.MaximumConnections)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	return helpers.SetJSON(d, "settings", res.Settings)
}

func resourceRACEndpointUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app, diags := resourceRACEndpointSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.RacApi.RacEndpointsUpdate(ctx, d.Id()).EndpointRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRACEndpointRead(ctx, d, m)
}

func resourceRACEndpointDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RacApi.RacEndpointsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
