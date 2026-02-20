package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Description:      helpers.EnumToDescription(api.AllowedPromptTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPromptTypeEnumEnumValues),
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"placeholder": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: helpers.DiffSuppressExpression,
			},
			"placeholder_expression": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"initial_value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: helpers.DiffSuppressExpression,
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
		Required:               new(d.Get("required").(bool)),
		PlaceholderExpression:  new(d.Get("placeholder_expression").(bool)),
		InitialValueExpression: new(d.Get("initial_value_expression").(bool)),
		SubText:                new(d.Get("sub_text").(string)),
		Placeholder:            helpers.GetP[string](d, "placeholder"),
		InitialValue:           helpers.GetP[string](d, "initial_value"),
		Order:                  helpers.GetIntP(d, "order"),
	}
	return &r
}

func resourceStagePromptFieldCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePromptFieldSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsCreate(ctx).PromptRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "field_key", res.FieldKey)
	helpers.SetWrapper(d, "label", res.Label)
	helpers.SetWrapper(d, "type", res.Type)
	helpers.SetWrapper(d, "required", res.Required)
	helpers.SetWrapper(d, "placeholder", res.Placeholder)
	helpers.SetWrapper(d, "placeholder_expression", res.PlaceholderExpression)
	helpers.SetWrapper(d, "initial_value", res.InitialValue)
	helpers.SetWrapper(d, "initial_value_expression", res.InitialValueExpression)
	helpers.SetWrapper(d, "sub_text", res.SubText)
	helpers.SetWrapper(d, "order", res.Order)
	return diags
}

func resourceStagePromptFieldUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePromptFieldSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPromptPromptsUpdate(ctx, d.Id()).PromptRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePromptFieldRead(ctx, d, m)
}

func resourceStagePromptFieldDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPromptPromptsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
