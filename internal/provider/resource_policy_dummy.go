package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
		ExecutionLogging: api.PtrBool(d.Get("execution_logging").(bool)),
		Result:           api.PtrBool(d.Get("result").(bool)),
	}

	if p, pSet := d.GetOk("wait_max"); pSet {
		r.WaitMax = api.PtrInt32(int32(p.(int)))
	}
	if p, pSet := d.GetOk("wait_min"); pSet {
		r.WaitMin = api.PtrInt32(int32(p.(int)))
	}
	return &r
}

func resourcePolicyDummyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyDummySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyCreate(ctx).DummyPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	setWrapper(d, "result", res.Result)
	setWrapper(d, "wait_min", res.WaitMin)
	setWrapper(d, "wait_max", res.WaitMax)
	return diags
}

func resourcePolicyDummyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyDummySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyUpdate(ctx, d.Id()).DummyPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesDummyDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
