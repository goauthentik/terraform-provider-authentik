package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
		Fields:             castSlice[string](d.Get("fields").([]interface{})),
		ValidationPolicies: castSlice[string](d.Get("validation_policies").([]interface{})),
	}
	return &r
}

func resourceStagePromptCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePromptSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptStagesCreate(ctx).PromptStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptStagesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	fields := castSlice[string](d.Get("fields").([]interface{}))
	setWrapper(d, "fields", listConsistentMerge(fields, res.Fields))
	validationPolicies := castSlice[string](d.Get("validation_policies").([]interface{}))
	setWrapper(d, "validation_policies", listConsistentMerge(validationPolicies, res.ValidationPolicies))
	return diags
}

func resourceStagePromptUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePromptSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptStagesUpdate(ctx, d.Id()).PromptStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptStagesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
