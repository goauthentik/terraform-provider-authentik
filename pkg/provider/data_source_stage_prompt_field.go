package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceStagePromptField() *schema.Resource {
	return &schema.Resource{
		Description: "Flows & Stages --- Get stage prompt fields by id or name",
		ReadContext: dataSourceStagePromptFieldRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"field_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"initial_value_expression": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"placeholder": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"placeholder_expression": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sub_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStagePromptFieldRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return diag.Errorf("Neither id nor name were provided")
	}

	var f api.Prompt

	if idOk {
		req := c.client.StagesAPI.StagesPromptPromptsRetrieve(ctx, id.(string))

		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		f = *res
	} else {
		req := c.client.StagesAPI.StagesPromptPromptsList(ctx)

		req = req.Name(name.(string))

		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		if len(res.Results) < 1 {
			return diag.Errorf("No matching stage prompt fields found")
		}

		if len(res.Results) > 1 {
			// In theory, impossible..
			return diag.Errorf("Multiple stage prompt fields found")
		}

		f = res.Results[0]
	}

	d.SetId(f.Pk)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetWrapper(d, "field_key", f.FieldKey)
	helpers.SetWrapper(d, "initial_value", f.InitialValue)
	helpers.SetWrapper(d, "initial_value_expression", f.InitialValueExpression)
	helpers.SetWrapper(d, "label", f.Label)
	helpers.SetWrapper(d, "order", f.Order)
	helpers.SetWrapper(d, "placeholder", f.Placeholder)
	helpers.SetWrapper(d, "placeholder_expression", f.PlaceholderExpression)
	helpers.SetWrapper(d, "required", f.Required)
	helpers.SetWrapper(d, "sub_text", f.SubText)
	helpers.SetWrapper(d, "type", f.Type)
	return diags
}
