package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStagePrompt() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStagePromptCreate,
		ReadContext:   resourceStagePromptRead,
		UpdateContext: resourceStagePromptUpdate,
		DeleteContext: resourceStagePromptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"validation_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStagePromptSchemaToProvider(d *schema.ResourceData) *api.PromptStageRequest {
	r := api.PromptStageRequest{
		Name:               d.Get("name").(string),
		Fields:             helpers.CastSlice[string](d, "fields"),
		ValidationPolicies: helpers.CastSlice[string](d, "validation_policies"),
	}
	return &r
}

func resourceStagePromptCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePromptSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptStagesCreate(ctx).PromptStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptStagesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	fields := helpers.CastSlice[string](d, "fields")
	helpers.SetWrapper(d, "fields", helpers.ListConsistentMerge(fields, res.Fields))
	validationPolicies := helpers.CastSlice[string](d, "validation_policies")
	helpers.SetWrapper(d, "validation_policies", helpers.ListConsistentMerge(validationPolicies, res.ValidationPolicies))
	return diags
}

func resourceStagePromptUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePromptSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptStagesUpdate(ctx, d.Id()).PromptStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptStagesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
