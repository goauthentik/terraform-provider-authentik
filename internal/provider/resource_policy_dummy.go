package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicyDummy() *schema.Resource {
	return &schema.Resource{
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

func resourcePolicyDummySchemaToProvider(d *schema.ResourceData) (*api.DummyPolicyRequest, diag.Diagnostics) {
	r := api.DummyPolicyRequest{
		ExecutionLogging: boolToPointer(d.Get("execution_logging").(bool)),
		Result:           boolToPointer(d.Get("result").(bool)),
	}
	r.Name.Set(stringToPointer(d.Get("name").(string)))

	if p, pSet := d.GetOk("wait_max"); pSet {
		r.WaitMax = intToPointer(p.(int))
	}
	if p, pSet := d.GetOk("wait_min"); pSet {
		r.WaitMin = intToPointer(p.(int))
	}

	return &r, nil
}

func resourcePolicyDummyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourcePolicyDummySchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PoliciesApi.PoliciesDummyCreate(ctx).DummyPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesDummyRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name.Get())
	d.Set("execution_logging", res.ExecutionLogging)
	d.Set("result", res.Result)
	d.Set("wait_min", res.WaitMin)
	d.Set("wait_max", res.WaitMax)
	return diags
}

func resourcePolicyDummyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourcePolicyDummySchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.PoliciesApi.PoliciesDummyUpdate(ctx, d.Id()).DummyPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyDummyRead(ctx, d, m)
}

func resourcePolicyDummyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesDummyDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
