package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStagePromptField() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStagePromptFieldCreate,
		ReadContext:   resourceStagePromptFieldRead,
		UpdateContext: resourceStagePromptFieldUpdate,
		DeleteContext: resourceStagePromptFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"field_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"placeholder": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceStagePromptFieldSchemaToProvider(d *schema.ResourceData) (*api.PromptRequest, diag.Diagnostics) {
	r := api.PromptRequest{
		FieldKey: d.Get("field_key").(string),
		Label:    d.Get("label").(string),
		Type:     api.PromptTypeEnum(d.Get("type").(string)),
		Required: boolToPointer(d.Get("required").(bool)),
	}

	if p, pSet := d.GetOk("placeholder"); pSet {
		r.Placeholder = stringToPointer(p.(string))
	}

	if o, oSet := d.GetOk("order"); oSet {
		r.Order = intToPointer(o.(int))
	}

	return &r, nil
}

func resourceStagePromptFieldCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStagePromptFieldSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesPromptPromptsCreate(ctx).PromptRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("field_key", res.FieldKey)
	d.Set("label", res.Label)
	d.Set("type", res.Type)
	d.Set("required", res.Required)
	d.Set("placeholder", res.Placeholder)
	d.Set("order", res.Order)
	return diags
}

func resourceStagePromptFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStagePromptFieldSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesPromptPromptsUpdate(ctx, d.Id()).PromptRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptPromptsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
