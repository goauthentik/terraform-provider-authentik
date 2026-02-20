package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageUserDelete() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageUserDeleteCreate,
		ReadContext:   resourceStageUserDeleteRead,
		UpdateContext: resourceStageUserDeleteUpdate,
		DeleteContext: resourceStageUserDeleteDelete,
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

func resourceStageUserDeleteSchemaToProvider(d *schema.ResourceData) *api.UserDeleteStageRequest {
	r := api.UserDeleteStageRequest{
		Name: d.Get("name").(string),
	}
	return &r
}

func resourceStageUserDeleteCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageUserDeleteSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserDeleteCreate(ctx).UserDeleteStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserDeleteRead(ctx, d, m)
}

func resourceStageUserDeleteRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserDeleteRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	return diags
}

func resourceStageUserDeleteUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageUserDeleteSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserDeleteUpdate(ctx, d.Id()).UserDeleteStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserDeleteRead(ctx, d, m)
}

func resourceStageUserDeleteDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserDeleteDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
