package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
		ExecutionLogging: new(d.Get("execution_logging").(bool)),
		CheckIp:          new(d.Get("check_ip").(bool)),
		CheckUsername:    new(d.Get("check_username").(bool)),
		Threshold:        helpers.GetIntP(d, "threshold"),
	}
	return &r
}

func resourcePolicyReputationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationCreate(ctx).ReputationPolicyRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "execution_logging", res.ExecutionLogging)
	helpers.SetWrapper(d, "check_ip", res.CheckIp)
	helpers.SetWrapper(d, "check_username", res.CheckUsername)
	helpers.SetWrapper(d, "threshold", res.Threshold)
	return diags
}

func resourcePolicyReputationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyReputationSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesReputationUpdate(ctx, d.Id()).ReputationPolicyRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyReputationRead(ctx, d, m)
}

func resourcePolicyReputationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesReputationDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
