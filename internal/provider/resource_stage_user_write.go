package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageUserWrite() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
			"user_creation_mode": {
				Type:     schema.TypeString,
				Default:  string(api.USERCREATIONMODEENUM_CREATE_WHEN_REQUIRED),
				Optional: true,
			},
			"create_users_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
		},
	}
}

func resourceStageUserWriteSchemaToProvider(d *schema.ResourceData) *api.UserWriteStageRequest {
	r := api.UserWriteStageRequest{
		Name:                  d.Get("name").(string),
		CreateUsersAsInactive: api.PtrBool(d.Get("create_users_as_inactive").(bool)),
		UserPathTemplate:      api.PtrString(d.Get("user_path_template").(string)),
		UserCreationMode:      api.UserCreationModeEnum(d.Get("user_creation_mode").(string)).Ptr(),
	}

	if h, hSet := d.GetOk("create_users_group"); hSet {
		r.CreateUsersGroup.Set(api.PtrString(h.(string)))
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

	setWrapper(d, "name", res.Name)
	setWrapper(d, "create_users_as_inactive", res.CreateUsersAsInactive)
	setWrapper(d, "create_users_group", res.CreateUsersGroup.Get())
	setWrapper(d, "user_path_template", res.GetUserPathTemplate())
	setWrapper(d, "user_creation_mode", res.GetUserCreationMode())
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
