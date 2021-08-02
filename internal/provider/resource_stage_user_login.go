package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStageUserLogin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageUserLoginCreate,
		ReadContext:   resourceStageUserLoginRead,
		UpdateContext: resourceStageUserLoginUpdate,
		DeleteContext: resourceStageUserLoginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"session_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "seconds=0",
			},
		},
	}
}

func resourceStageUserLoginSchemaToProvider(d *schema.ResourceData) (*api.UserLoginStageRequest, diag.Diagnostics) {
	r := api.UserLoginStageRequest{
		Name:            d.Get("name").(string),
		SessionDuration: stringToPointer(d.Get("session_duration").(string)),
	}

	return &r, nil
}

func resourceStageUserLoginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageUserLoginSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesUserLoginCreate(ctx).UserLoginStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLoginRead(ctx, d, m)
}

func resourceStageUserLoginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserLoginRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("session_duration", res.SessionDuration)
	return diags
}

func resourceStageUserLoginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStageUserLoginSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesUserLoginUpdate(ctx, d.Id()).UserLoginStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLoginRead(ctx, d, m)
}

func resourceStageUserLoginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserLoginDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
