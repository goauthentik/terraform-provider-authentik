package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStagePromptField() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStagePromptFieldCreate,
		ReadContext:   resourceStagePromptFieldRead,
		UpdateContext: resourceStagePromptFieldUpdate,
		DeleteContext: resourceStagePromptFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      EnumToDescription(api.AllowedPromptTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPromptTypeEnumEnumValues),
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"placeholder": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
			"placeholder_expression": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"initial_value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressExpression,
			},
			"initial_value_expression": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sub_text": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceStagePromptFieldSchemaToProvider(d *schema.ResourceData) *api.PromptRequest {
	r := api.PromptRequest{
		Name:                   d.Get("name").(string),
		FieldKey:               d.Get("field_key").(string),
		Label:                  d.Get("label").(string),
		Type:                   api.PromptTypeEnum(d.Get("type").(string)),
		Required:               api.PtrBool(d.Get("required").(bool)),
		PlaceholderExpression:  api.PtrBool(d.Get("placeholder_expression").(bool)),
		InitialValueExpression: api.PtrBool(d.Get("initial_value_expression").(bool)),
		SubText:                api.PtrString(d.Get("sub_text").(string)),
	}

	if p, pSet := d.GetOk("placeholder"); pSet {
		r.Placeholder = api.PtrString(p.(string))
	}

	if p, pSet := d.GetOk("initial_value"); pSet {
		r.InitialValue = api.PtrString(p.(string))
	}

	if o, oSet := d.GetOk("order"); oSet {
		r.Order = api.PtrInt32(int32(o.(int)))
	}
	return &r
}

func resourceStagePromptFieldCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePromptFieldSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsCreate(ctx).PromptRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "field_key", res.FieldKey)
	setWrapper(d, "label", res.Label)
	setWrapper(d, "type", res.Type)
	setWrapper(d, "required", res.Required)
	setWrapper(d, "placeholder", res.Placeholder)
	setWrapper(d, "placeholder_expression", res.PlaceholderExpression)
	setWrapper(d, "initial_value", res.InitialValue)
	setWrapper(d, "initial_value_expression", res.InitialValueExpression)
	setWrapper(d, "sub_text", res.SubText)
	setWrapper(d, "order", res.Order)
	return diags
}

func resourceStagePromptFieldUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePromptFieldSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsUpdate(ctx, d.Id()).PromptRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptPromptsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
