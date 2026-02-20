package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePolicyDummy() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyDummyCreate,
		ReadContext:   resourcePolicyDummyRead,
		UpdateContext: resourcePolicyDummyUpdate,
		DeleteContext: resourcePolicyDummyDelete,
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
			"result": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"wait_min": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"wait_max": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
	}
}

func resourcePolicyDummySchemaToProvider(d *schema.ResourceData) *api.DummyPolicyRequest {
	r := api.DummyPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: new(d.Get("execution_logging").(bool)),
		Result:           new(d.Get("result").(bool)),
		WaitMin:          helpers.GetIntP(d, "wait_min"),
		WaitMax:          helpers.GetIntP(d, "wait_max"),
	}
	return &r
}

func resourcePolicyDummyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyDummySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyCreate(ctx).DummyPolicyRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "execution_logging", res.ExecutionLogging)
	helpers.SetWrapper(d, "result", res.Result)
	helpers.SetWrapper(d, "wait_min", res.WaitMin)
	helpers.SetWrapper(d, "wait_max", res.WaitMax)
	return diags
}

func resourcePolicyDummyUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyDummySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyUpdate(ctx, d.Id()).DummyPolicyRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesDummyDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
