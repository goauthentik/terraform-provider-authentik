package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorEndpointGDTC() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorEndpointGDTCCreate,
		ReadContext:   resourceStageAuthenticatorEndpointGDTCRead,
		UpdateContext: resourceStageAuthenticatorEndpointGDTCUpdate,
		DeleteContext: resourceStageAuthenticatorEndpointGDTCDelete,
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
			"credentials": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d *schema.ResourceData) (*api.AuthenticatorEndpointGDTCStageRequest, diag.Diagnostics) {
	r := api.AuthenticatorEndpointGDTCStageRequest{
		Name: d.Get("name").(string),
	}
	attr, err := helpers.GetJSON[map[string]any](d, ("credentials"))
	r.Credentials = attr
	return &r, err
}

func resourceStageAuthenticatorEndpointGDTCCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcCreate(ctx).AuthenticatorEndpointGDTCStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEndpointGDTCRead(ctx, d, m)
}

func resourceStageAuthenticatorEndpointGDTCRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	return helpers.SetJSON(d, "credentials", res.Credentials)
}

func resourceStageAuthenticatorEndpointGDTCUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageAuthenticatorEndpointGDTCSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcUpdate(ctx, d.Id()).AuthenticatorEndpointGDTCStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEndpointGDTCRead(ctx, d, m)
}

func resourceStageAuthenticatorEndpointGDTCDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorEndpointGdtcDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
