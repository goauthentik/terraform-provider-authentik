package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourcePolicyReputation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyReputationCreate,
		ReadContext:   resourcePolicyReputationRead,
		UpdateContext: resourcePolicyReputationUpdate,
		DeleteContext: resourcePolicyReputationDelete,
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
			"check_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"check_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
		},
	}
}

func resourcePolicyReputationSchemaToProvider(d *schema.ResourceData) *api.ReputationPolicyRequest {
	r := api.ReputationPolicyRequest{
		ExecutionLogging: boolToPointer(d.Get("execution_logging").(bool)),
		CheckIp:          boolToPointer(d.Get("check_ip").(bool)),
		CheckUsername:    boolToPointer(d.Get("check_username").(bool)),
	}
	r.Name.Set(stringToPointer(d.Get("name").(string)))

	if p, pSet := d.GetOk("threshold"); pSet {
		r.Threshold = intToPointer(p.(int))
	}

	return &r
}

func resourcePolicyReputationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationCreate(ctx).ReputationPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name.Get())
	d.Set("execution_logging", res.ExecutionLogging)
	d.Set("check_ip", res.CheckIp)
	d.Set("check_username", res.CheckUsername)
	d.Set("threshold", res.Threshold)
	return diags
}

func resourcePolicyReputationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationUpdate(ctx, d.Id()).ReputationPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesReputationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
