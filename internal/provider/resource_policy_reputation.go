package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyReputation() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
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
		Name:             d.Get("name").(string),
		ExecutionLogging: api.PtrBool(d.Get("execution_logging").(bool)),
		CheckIp:          api.PtrBool(d.Get("check_ip").(bool)),
		CheckUsername:    api.PtrBool(d.Get("check_username").(bool)),
	}

	if p, pSet := d.GetOk("threshold"); pSet {
		r.Threshold = api.PtrInt32(int32(p.(int)))
	}
	return &r
}

func resourcePolicyReputationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationCreate(ctx).ReputationPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	setWrapper(d, "check_ip", res.CheckIp)
	setWrapper(d, "check_username", res.CheckUsername)
	setWrapper(d, "threshold", res.Threshold)
	return diags
}

func resourcePolicyReputationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationUpdate(ctx, d.Id()).ReputationPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesReputationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
