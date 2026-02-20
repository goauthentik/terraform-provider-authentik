package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePolicyExpression() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyExpressionCreate,
		ReadContext:   resourcePolicyExpressionRead,
		UpdateContext: resourcePolicyExpressionUpdate,
		DeleteContext: resourcePolicyExpressionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"expression": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: helpers.DiffSuppressExpression,
			},
		},
	}
}

func resourcePolicyExpressionSchemaToProvider(d *schema.ResourceData) *api.ExpressionPolicyRequest {
	r := api.ExpressionPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: new(d.Get("execution_logging").(bool)),
		Expression:       d.Get("expression").(string),
	}
	return &r
}

func resourcePolicyExpressionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyExpressionSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesExpressionCreate(ctx).ExpressionPolicyRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyExpressionRead(ctx, d, m)
}

func resourcePolicyExpressionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesExpressionRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "execution_logging", res.ExecutionLogging)
	helpers.SetWrapper(d, "expression", res.Expression)
	return diags
}

func resourcePolicyExpressionUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyExpressionSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesExpressionUpdate(ctx, d.Id()).ExpressionPolicyRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyExpressionRead(ctx, d, m)
}

func resourcePolicyExpressionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesExpressionDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
