package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStagePrompt() *schema.Resource {
	return &schema.Resource{
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

func resourceStagePromptSchemaToProvider(d *schema.ResourceData) (*api.PromptStageRequest, diag.Diagnostics) {
	r := api.PromptStageRequest{
		Name: d.Get("name").(string),
	}

	fields := make([]string, 0)
	for _, fieldsS := range d.Get("fields").([]interface{}) {
		fields = append(fields, fieldsS.(string))
	}
	r.Fields = fields

	vp := make([]string, 0)
	for _, vpS := range d.Get("validation_policies").([]interface{}) {
		vp = append(vp, vpS.(string))
	}
	r.ValidationPolicies = &vp

	return &r, nil
}

func resourceStagePromptCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStagePromptSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesPromptStagesCreate(ctx).PromptStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptStagesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("fields", res.Fields)
	d.Set("validation_policies", res.ValidationPolicies)
	return diags
}

func resourceStagePromptUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStagePromptSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesPromptStagesUpdate(ctx, d.Id()).PromptStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptRead(ctx, d, m)
}

func resourceStagePromptDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptStagesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
