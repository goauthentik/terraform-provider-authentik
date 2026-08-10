package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAccountLockdown() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAccountLockdownCreate,
		ReadContext:   resourceStageAccountLockdownRead,
		UpdateContext: resourceStageAccountLockdownUpdate,
		DeleteContext: resourceStageAccountLockdownDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deactivate_user": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"set_unusable_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"delete_sessions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"revoke_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"self_service_completion_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageAccountLockdownSchemaToProvider(d *schema.ResourceData) *api.AccountLockdownStageRequest {
	r := api.AccountLockdownStageRequest{
		Name:                      d.Get("name").(string),
		DeactivateUser:            new(d.Get("deactivate_user").(bool)),
		SetUnusablePassword:       new(d.Get("set_unusable_password").(bool)),
		DeleteSessions:            new(d.Get("delete_sessions").(bool)),
		RevokeTokens:              new(d.Get("revoke_tokens").(bool)),
		SelfServiceCompletionFlow: *api.NewNullableString(helpers.GetP[string](d, "self_service_completion_flow")),
	}
	return &r
}

func resourceStageAccountLockdownCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAccountLockdownSchemaToProvider(d)

	res, hr, err := c.client.StagesAPI.StagesAccountLockdownCreate(ctx).AccountLockdownStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAccountLockdownRead(ctx, d, m)
}

func resourceStageAccountLockdownRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesAPI.StagesAccountLockdownRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "deactivate_user", res.DeactivateUser)
	helpers.SetWrapper(d, "set_unusable_password", res.SetUnusablePassword)
	helpers.SetWrapper(d, "delete_sessions", res.DeleteSessions)
	helpers.SetWrapper(d, "revoke_tokens", res.RevokeTokens)
	helpers.SetWrapper(d, "self_service_completion_flow", res.SelfServiceCompletionFlow.Get())
	return diags
}

func resourceStageAccountLockdownUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAccountLockdownSchemaToProvider(d)

	res, hr, err := c.client.StagesAPI.StagesAccountLockdownUpdate(ctx, d.Id()).AccountLockdownStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAccountLockdownRead(ctx, d, m)
}

func resourceStageAccountLockdownDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesAPI.StagesAccountLockdownDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
