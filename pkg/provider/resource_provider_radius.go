package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceProviderRadius() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderRadiusCreate,
		ReadContext:   resourceProviderRadiusRead,
		UpdateContext: resourceProviderRadiusUpdate,
		DeleteContext: resourceProviderRadiusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"invalidation_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"client_networks": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0/0, ::/0",
			},
			"shared_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"mfa_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProviderRadiusSchemaToProvider(d *schema.ResourceData) *api.RadiusProviderRequest {
	r := api.RadiusProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("authorization_flow").(string),
		InvalidationFlow:  d.Get("invalidation_flow").(string),
		ClientNetworks:    new(d.Get("client_networks").(string)),
		SharedSecret:      new(d.Get("shared_secret").(string)),
		MfaSupport:        new(d.Get("mfa_support").(bool)),
		PropertyMappings:  helpers.CastSlice[string](d, "property_mappings"),
		Certificate:       *api.NewNullableString(helpers.GetP[string](d, "certificate")),
	}
	return &r
}

func resourceProviderRadiusCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersRadiusCreate(ctx).RadiusProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRadiusRead(ctx, d, m)
}

func resourceProviderRadiusRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersRadiusRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "invalidation_flow", res.InvalidationFlow)
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	helpers.SetWrapper(d, "client_networks", res.ClientNetworks)
	helpers.SetWrapper(d, "shared_secret", res.SharedSecret)
	helpers.SetWrapper(d, "mfa_support", res.MfaSupport)
	helpers.SetWrapper(d, "certificate", res.Certificate.Get())
	return diags
}

func resourceProviderRadiusUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersRadiusUpdate(ctx, int32(id)).RadiusProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRadiusRead(ctx, d, m)
}

func resourceProviderRadiusDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersRadiusDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
