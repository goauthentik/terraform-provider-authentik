package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStageUserLogout() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageUserLogoutCreate,
		ReadContext:   resourceStageUserLogoutRead,
		UpdateContext: resourceStageUserLogoutUpdate,
		DeleteContext: resourceStageUserLogoutDelete,
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

func resourceStageUserLogoutSchemaToProvider(d *schema.ResourceData) (*api.UserLogoutStageRequest, diag.Diagnostics) {
	r := api.UserLogoutStageRequest{
		Name: d.Get("name").(string),
	}

	return &r, nil
}

func resourceStageUserLogoutCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageUserLogoutSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesUserLogoutCreate(ctx).UserLogoutStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStageUserLogoutRead(ctx, d, m)
}

func resourceStageUserLogoutRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserLogoutRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	return diags
}

func resourceStageUserLogoutUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStageUserLogoutSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesUserLogoutUpdate(ctx, d.Id()).UserLogoutStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStageUserLogoutRead(ctx, d, m)
}

func resourceStageUserLogoutDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserLogoutDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
