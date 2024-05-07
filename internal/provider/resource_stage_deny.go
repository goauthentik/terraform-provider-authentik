package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageDeny() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageDenyCreate,
		ReadContext:   resourceStageDenyRead,
		UpdateContext: resourceStageDenyUpdate,
		DeleteContext: resourceStageDenyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deny_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageDenySchemaToProvider(d *schema.ResourceData) *api.DenyStageRequest {
	r := api.DenyStageRequest{
		Name: d.Get("name").(string),
	}
	if m, mok := d.GetOk("deny_message"); mok {
		r.DenyMessage = api.PtrString(m.(string))
	}
	return &r
}

func resourceStageDenyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageDenySchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesDenyCreate(ctx).DenyStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageDenyRead(ctx, d, m)
}

func resourceStageDenyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesDenyRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "deny_message", res.GetDenyMessage())
	return diags
}

func resourceStageDenyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageDenySchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesDenyUpdate(ctx, d.Id()).DenyStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageDenyRead(ctx, d, m)
}

func resourceStageDenyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesDenyDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
