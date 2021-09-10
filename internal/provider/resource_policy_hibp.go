package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourcePolicyHaveIBeenPwend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyHaveIBeenPwendCreate,
		ReadContext:   resourcePolicyHaveIBeenPwendRead,
		UpdateContext: resourcePolicyHaveIBeenPwendUpdate,
		DeleteContext: resourcePolicyHaveIBeenPwendDelete,
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
			"password_field": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password",
			},
			"allowed_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourcePolicyHaveIBeenPwendSchemaToProvider(d *schema.ResourceData) *api.HaveIBeenPwendPolicyRequest {
	r := api.HaveIBeenPwendPolicyRequest{
		ExecutionLogging: boolToPointer(d.Get("execution_logging").(bool)),
	}
	r.Name.Set(stringToPointer(d.Get("name").(string)))

	if p, sSet := d.GetOk("allowed_count"); sSet {
		r.AllowedCount = intToPointer(p.(int))
	}
	if s, sSet := d.GetOk("password_field"); sSet {
		r.PasswordField = stringToPointer(s.(string))
	}
	return &r
}

func resourcePolicyHaveIBeenPwendCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyHaveIBeenPwendSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesHaveibeenpwnedCreate(ctx).HaveIBeenPwendPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyHaveIBeenPwendRead(ctx, d, m)
}

func resourcePolicyHaveIBeenPwendRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesHaveibeenpwnedRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name.Get())
	d.Set("execution_logging", res.ExecutionLogging)
	d.Set("password_field", res.PasswordField)
	d.Set("allowed_count", res.AllowedCount)
	return diags
}

func resourcePolicyHaveIBeenPwendUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyHaveIBeenPwendSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesHaveibeenpwnedUpdate(ctx, d.Id()).HaveIBeenPwendPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyHaveIBeenPwendRead(ctx, d, m)
}

func resourcePolicyHaveIBeenPwendDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesHaveibeenpwnedDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
