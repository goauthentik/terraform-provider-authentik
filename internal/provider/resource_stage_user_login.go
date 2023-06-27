package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
			"terminate_other_sessions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageUserLoginSchemaToProvider(d *schema.ResourceData) *api.UserLoginStageRequest {
	r := api.UserLoginStageRequest{
		Name:                   d.Get("name").(string),
		SessionDuration:        api.PtrString(d.Get("session_duration").(string)),
		TerminateOtherSessions: api.PtrBool(d.Get("terminate_other_sessions").(bool)),
	}
	return &r
}

func resourceStageUserLoginCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageUserLoginSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserLoginCreate(ctx).UserLoginStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLoginRead(ctx, d, m)
}

func resourceStageUserLoginRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesUserLoginRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "session_duration", res.SessionDuration)
	setWrapper(d, "terminate_other_sessions", res.TerminateOtherSessions)
	return diags
}

func resourceStageUserLoginUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageUserLoginSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesUserLoginUpdate(ctx, d.Id()).UserLoginStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageUserLoginRead(ctx, d, m)
}

func resourceStageUserLoginDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesUserLoginDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
