package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStageAuthenticatorWebAuthn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageAuthenticatorWebAuthnCreate,
		ReadContext:   resourceStageAuthenticatorWebAuthnRead,
		UpdateContext: resourceStageAuthenticatorWebAuthnUpdate,
		DeleteContext: resourceStageAuthenticatorWebAuthnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageAuthenticatorWebAuthnSchemaToProvider(d *schema.ResourceData) (*api.AuthenticateWebAuthnStageRequest, diag.Diagnostics) {
	r := api.AuthenticateWebAuthnStageRequest{
		Name: d.Get("name").(string),
	}

	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(stringToPointer(h.(string)))
	}
	return &r, nil
}

func resourceStageAuthenticatorWebAuthnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnCreate(ctx).AuthenticateWebAuthnStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	if res.ConfigureFlow.IsSet() {
		d.Set("configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorWebAuthnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnUpdate(ctx, d.Id()).AuthenticateWebAuthnStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
