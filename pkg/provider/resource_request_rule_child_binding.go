package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRequestRuleChildBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceRequestRuleChildBindingCreate,
		ReadContext:   resourceRequestRuleChildBindingRead,
		UpdateContext: resourceRequestRuleChildBindingUpdate,
		DeleteContext: resourceRequestRuleChildBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"binding": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the authentik_request_rule_binding this child binding extends.",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the additional object this child binding makes requestable alongside the parent binding's target.",
			},
		},
	}
}

func resourceRequestRuleChildBindingSchemaToModel(d *schema.ResourceData) *api.RequestRuleChildBindingRequest {
	m := api.RequestRuleChildBindingRequest{
		Binding: d.Get("binding").(string),
		Target:  d.Get("target").(string),
	}
	return &m
}

func resourceRequestRuleChildBindingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleChildBindingSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRuleChildBindingsCreate(ctx).RequestRuleChildBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleChildBindingRead(ctx, d, m)
}

func resourceRequestRuleChildBindingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RequestsAPI.RequestsRuleChildBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "binding", res.Binding)
	helpers.SetWrapper(d, "target", res.Target)
	return diags
}

func resourceRequestRuleChildBindingUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleChildBindingSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRuleChildBindingsUpdate(ctx, d.Id()).RequestRuleChildBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleChildBindingRead(ctx, d, m)
}

func resourceRequestRuleChildBindingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RequestsAPI.RequestsRuleChildBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
