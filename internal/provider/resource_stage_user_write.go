package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func resourceStageUserWrite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageUserWriteCreate,
		ReadContext:   resourceStageUserWriteRead,
		UpdateContext: resourceStageUserWriteUpdate,
		DeleteContext: resourceStageUserWriteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_users_as_inactive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"create_users_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageUserWriteSchemaToProvider(d *schema.ResourceData) *api.UserWriteStageRequest {
	r := api.UserWriteStageRequest{
		Name:                  d.Get("name").(string),
		CreateUsersAsInactive: boolToPointer(d.Get("create_users_as_inactive").(bool)),
	}

	if h, hSet := d.GetOk("create_users_group"); hSet {
		r.CreateUsersGroup.Set(stringToPointer(h.(string)))
	}
	return &r
}

func resourceStageUserWriteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageUserWriteSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserWriteCreate(ctx).UserWriteStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserWriteRead(ctx, d, m)
}

func resourceStageUserWriteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserWriteRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("create_users_as_inactive", res.CreateUsersAsInactive)
	d.Set("create_users_group", res.CreateUsersGroup)
	return diags
}

func resourceStageUserWriteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageUserWriteSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserWriteUpdate(ctx, d.Id()).UserWriteStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserWriteRead(ctx, d, m)
}

func resourceStageUserWriteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserWriteDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
