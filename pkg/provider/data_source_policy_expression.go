package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourcePolicyExpression() *schema.Resource {
	return &schema.Resource{
		Description: "Customization --- Get policy expressions by id or name",
		ReadContext: dataSourcePolicyExpressionRead,
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
			"bound_to": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"component": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_model_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verbose_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verbose_name_plural": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyExpressionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return diag.Errorf("Neither id nor name were provided")
	}

	req := c.client.PoliciesAPI.PoliciesExpressionList(ctx)

	if idOk {
		req = req.PolicyUuid(id.(string))
	} else {
		req = req.Name(name.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching expression policies found")
	}

	if len(res.Results) > 1 {
		// In theory, impossible..
		return diag.Errorf("Multiple expression policies found")
	}

	f := res.Results[0]
	d.SetId(f.Pk)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetWrapper(d, "bound_to", f.BoundTo)
	helpers.SetWrapper(d, "component", f.Component)
	helpers.SetWrapper(d, "execution_logging", f.ExecutionLogging)
	helpers.SetWrapper(d, "expression", f.Expression)
	helpers.SetWrapper(d, "meta_model_name", f.MetaModelName)
	helpers.SetWrapper(d, "verbose_name", f.VerboseName)
	helpers.SetWrapper(d, "verbose_name_plural", f.VerboseNamePlural)
	return diags
}
