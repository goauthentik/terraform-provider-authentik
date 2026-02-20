package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageInvitation() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageInvitationCreate,
		ReadContext:   resourceStageInvitationRead,
		UpdateContext: resourceStageInvitationUpdate,
		DeleteContext: resourceStageInvitationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"continue_flow_without_invitation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageInvitationSchemaToProvider(d *schema.ResourceData) *api.InvitationStageRequest {
	r := api.InvitationStageRequest{
		Name:                          d.Get("name").(string),
		ContinueFlowWithoutInvitation: new(d.Get("continue_flow_without_invitation").(bool)),
	}
	return &r
}

func resourceStageInvitationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageInvitationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesInvitationStagesCreate(ctx).InvitationStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageInvitationRead(ctx, d, m)
}

func resourceStageInvitationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesInvitationStagesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "continue_flow_without_invitation", res.ContinueFlowWithoutInvitation)
	return diags
}

func resourceStageInvitationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageInvitationSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesInvitationStagesUpdate(ctx, d.Id()).InvitationStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageInvitationRead(ctx, d, m)
}

func resourceStageInvitationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesInvitationStagesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
