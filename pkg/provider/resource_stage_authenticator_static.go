package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorStatic() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorStaticCreate,
		ReadContext:   resourceStageAuthenticatorStaticRead,
		UpdateContext: resourceStageAuthenticatorStaticUpdate,
		DeleteContext: resourceStageAuthenticatorStaticDelete,
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
			"token_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6,
			},
			"token_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  12,
			},
		},
	}
}

func resourceStageAuthenticatorStaticSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorStaticStageRequest {
	r := api.AuthenticatorStaticStageRequest{
		Name:          d.Get("name").(string),
		TokenCount:    new(int32(d.Get("token_count").(int))),
		TokenLength:   new(int32(d.Get("token_length").(int))),
		FriendlyName:  helpers.GetP[string](d, "friendly_name"),
		ConfigureFlow: *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),
	}
	return &r
}

func resourceStageAuthenticatorStaticCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorStaticSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorStaticCreate(ctx).AuthenticatorStaticStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorStaticRead(ctx, d, m)
}

func resourceStageAuthenticatorStaticRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorStaticRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "token_count", res.TokenCount)
	helpers.SetWrapper(d, "token_length", res.TokenLength)
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	return diags
}

func resourceStageAuthenticatorStaticUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorStaticSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorStaticUpdate(ctx, d.Id()).AuthenticatorStaticStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorStaticRead(ctx, d, m)
}

func resourceStageAuthenticatorStaticDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorStaticDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
