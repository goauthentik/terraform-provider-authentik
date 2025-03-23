package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyExpiry() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyExpiryCreate,
		ReadContext:   resourcePolicyExpiryRead,
		UpdateContext: resourcePolicyExpiryUpdate,
		DeleteContext: resourcePolicyExpiryDelete,
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
			"days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"deny_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourcePolicyExpirySchemaToProvider(d *schema.ResourceData) *api.PasswordExpiryPolicyRequest {
	r := api.PasswordExpiryPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: api.PtrBool(d.Get("execution_logging").(bool)),
		Days:             int32(d.Get("days").(int)),
		DenyOnly:         api.PtrBool(d.Get("deny_only").(bool)),
	}
	return &r
}

func resourcePolicyExpiryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyExpirySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordExpiryCreate(ctx).PasswordExpiryPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyExpiryRead(ctx, d, m)
}

func resourcePolicyExpiryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordExpiryRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	setWrapper(d, "days", res.Days)
	setWrapper(d, "deny_only", res.DenyOnly)
	return diags
}

func resourcePolicyExpiryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyExpirySchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordExpiryUpdate(ctx, d.Id()).PasswordExpiryPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyExpiryRead(ctx, d, m)
}

func resourcePolicyExpiryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesPasswordExpiryDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
