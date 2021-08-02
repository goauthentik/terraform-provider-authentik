package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicyPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPasswordCreate,
		ReadContext:   resourcePolicyPasswordRead,
		UpdateContext: resourcePolicyPasswordUpdate,
		DeleteContext: resourcePolicyPasswordDelete,
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
			"error_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"amount_uppercase": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"amount_lowercase": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"amount_symbols": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"length_min": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"symbol_charset": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "!\\\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
			},
		},
	}
}

func resourcePolicyPasswordSchemaToProvider(d *schema.ResourceData) (*api.PasswordPolicyRequest, diag.Diagnostics) {
	r := api.PasswordPolicyRequest{
		ExecutionLogging: boolToPointer(d.Get("execution_logging").(bool)),
	}
	r.Name.Set(stringToPointer(d.Get("name").(string)))

	if s, sSet := d.GetOk("symbol_charset"); sSet {
		r.SymbolCharset = stringToPointer(s.(string))
	}
	if s, sSet := d.GetOk("password_field"); sSet {
		r.PasswordField = stringToPointer(s.(string))
	}
	if s, sSet := d.GetOk("error_message"); sSet {
		r.ErrorMessage = s.(string)
	}

	if p, pSet := d.GetOk("amount_uppercase"); pSet {
		r.AmountUppercase = intToPointer(p.(int))
	}
	if p, pSet := d.GetOk("amount_lowercase"); pSet {
		r.AmountLowercase = intToPointer(p.(int))
	}
	if p, pSet := d.GetOk("amount_symbols"); pSet {
		r.AmountSymbols = intToPointer(p.(int))
	}
	if p, pSet := d.GetOk("length_min"); pSet {
		r.LengthMin = intToPointer(p.(int))
	}

	return &r, nil
}

func resourcePolicyPasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourcePolicyPasswordSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordCreate(ctx).PasswordPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyPasswordRead(ctx, d, m)
}

func resourcePolicyPasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name.Get())
	d.Set("execution_logging", res.ExecutionLogging)
	d.Set("password_field", res.PasswordField)
	d.Set("error_message", res.ErrorMessage)
	d.Set("amount_uppercase", res.AmountUppercase)
	d.Set("amount_lowercase", res.AmountLowercase)
	d.Set("amount_symbols", res.AmountSymbols)
	d.Set("length_min", res.LengthMin)
	d.Set("symbol_charset", res.SymbolCharset)
	return diags
}

func resourcePolicyPasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourcePolicyPasswordSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordUpdate(ctx, d.Id()).PasswordPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyPasswordRead(ctx, d, m)
}

func resourcePolicyPasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesPasswordDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
