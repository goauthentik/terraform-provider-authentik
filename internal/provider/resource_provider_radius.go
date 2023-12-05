package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
		},
	}
}

func resourceProviderRadiusSchemaToProvider(d *schema.ResourceData) *api.RadiusProviderRequest {
	r := api.RadiusProviderRequest{
		Name:              d.Get("name").(string),
		AuthorizationFlow: d.Get("authorization_flow").(string),
		ClientNetworks:    api.PtrString(d.Get("client_networks").(string)),
		SharedSecret:      api.PtrString(d.Get("shared_secret").(string)),
		MfaSupport:        api.PtrBool(d.Get("mfa_support").(bool)),
	}
	return &r
}

func resourceProviderRadiusCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersRadiusCreate(ctx).RadiusProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRadiusRead(ctx, d, m)
}

func resourceProviderRadiusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersRadiusRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "authorization_flow", res.AuthorizationFlow)
	setWrapper(d, "client_networks", res.ClientNetworks)
	setWrapper(d, "shared_secret", res.SharedSecret)
	setWrapper(d, "mfa_support", res.MfaSupport)
	return diags
}

func resourceProviderRadiusUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	app := resourceProviderRadiusSchemaToProvider(d)

	res, hr, err := c.client.ProvidersApi.ProvidersRadiusUpdate(ctx, int32(id)).RadiusProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRadiusRead(ctx, d, m)
}

func resourceProviderRadiusDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersRadiusDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
