package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyUniquePassword() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyUniquePasswordCreate,
		ReadContext:   resourcePolicyUniquePasswordRead,
		UpdateContext: resourcePolicyUniquePasswordUpdate,
		DeleteContext: resourcePolicyUniquePasswordDelete,
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
			"num_historical_passwords": {
				Type:     schema.TypeInt,
				Default:  1,
				Optional: true,
			},
		},
	}
}

func resourcePolicyUniquePasswordSchemaToProvider(d *schema.ResourceData) *api.UniquePasswordPolicyRequest {
	r := api.UniquePasswordPolicyRequest{
		Name:                   d.Get("name").(string),
		ExecutionLogging:       api.PtrBool(d.Get("execution_logging").(bool)),
		PasswordField:          api.PtrString(d.Get("password_field").(string)),
		NumHistoricalPasswords: api.PtrInt32(int32(d.Get("num_historical_passwords").(int))),
	}
	return &r
}

func resourcePolicyUniquePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyUniquePasswordSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesUniquePasswordCreate(ctx).UniquePasswordPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyUniquePasswordRead(ctx, d, m)
}

func resourcePolicyUniquePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesUniquePasswordRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	setWrapper(d, "password_field", res.PasswordField)
	setWrapper(d, "num_historical_passwords", res.NumHistoricalPasswords)
	return diags
}

func resourcePolicyUniquePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyUniquePasswordSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesUniquePasswordUpdate(ctx, d.Id()).UniquePasswordPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyUniquePasswordRead(ctx, d, m)
}

func resourcePolicyUniquePasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesUniquePasswordDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
