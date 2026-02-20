package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageUserLogout() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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

func resourceStageUserLogoutSchemaToProvider(d *schema.ResourceData) *api.UserLogoutStageRequest {
	r := api.UserLogoutStageRequest{
		Name: d.Get("name").(string),
	}
	return &r
}

func resourceStageUserLogoutCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageUserLogoutSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserLogoutCreate(ctx).UserLogoutStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLogoutRead(ctx, d, m)
}

func resourceStageUserLogoutRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserLogoutRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	return diags
}

func resourceStageUserLogoutUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageUserLogoutSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserLogoutUpdate(ctx, d.Id()).UserLogoutStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLogoutRead(ctx, d, m)
}

func resourceStageUserLogoutDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserLogoutDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
