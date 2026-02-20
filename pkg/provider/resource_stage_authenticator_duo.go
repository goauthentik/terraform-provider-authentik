package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorDuo() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorDuoCreate,
		ReadContext:   resourceStageAuthenticatorDuoRead,
		UpdateContext: resourceStageAuthenticatorDuoUpdate,
		DeleteContext: resourceStageAuthenticatorDuoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"admin_integration_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"api_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceStageAuthenticatorDuoSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorDuoStageRequest {
	r := api.AuthenticatorDuoStageRequest{
		Name:                d.Get("name").(string),
		ClientId:            d.Get("client_id").(string),
		ClientSecret:        d.Get("client_secret").(string),
		ApiHostname:         d.Get("api_hostname").(string),
		FriendlyName:        helpers.GetP[string](d, "friendly_name"),
		AdminIntegrationKey: helpers.GetP[string](d, "admin_integration_key"),
		AdminSecretKey:      helpers.GetP[string](d, "admin_secret_key"),
		ConfigureFlow:       *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),
	}
	return &r
}

func resourceStageAuthenticatorDuoCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorDuoSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoCreate(ctx).AuthenticatorDuoStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorDuoRead(ctx, d, m)
}

func resourceStageAuthenticatorDuoRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "client_id", res.ClientId)
	helpers.SetWrapper(d, "admin_integration_key", res.AdminIntegrationKey)
	helpers.SetWrapper(d, "api_hostname", res.ApiHostname)
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	return diags
}

func resourceStageAuthenticatorDuoUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorDuoSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorDuoUpdate(ctx, d.Id()).AuthenticatorDuoStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorDuoRead(ctx, d, m)
}

func resourceStageAuthenticatorDuoDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorDuoDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
