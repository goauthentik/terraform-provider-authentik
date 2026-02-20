package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePolicyPassword() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
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

			"check_static_rules": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"check_have_i_been_pwned": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"check_zxcvbn": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
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
			"amount_digits": {
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

			"hibp_allowed_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"zxcvbn_score_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}
}

func resourcePolicyPasswordSchemaToProvider(d *schema.ResourceData) *api.PasswordPolicyRequest {
	r := api.PasswordPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: new(d.Get("execution_logging").(bool)),
		PasswordField:    helpers.GetP[string](d, "password_field"),

		CheckStaticRules:    helpers.GetP[bool](d, "check_static_rules"),
		CheckHaveIBeenPwned: helpers.GetP[bool](d, "check_have_i_been_pwned"),
		CheckZxcvbn:         helpers.GetP[bool](d, "check_zxcvbn"),

		SymbolCharset:   helpers.GetP[string](d, "symbol_charset"),
		ErrorMessage:    helpers.GetP[string](d, "error_message"),
		AmountUppercase: helpers.GetIntP(d, "amount_uppercase"),
		AmountDigits:    helpers.GetIntP(d, "amount_digits"),
		AmountLowercase: helpers.GetIntP(d, "amount_lowercase"),
		AmountSymbols:   helpers.GetIntP(d, "amount_symbols"),
		LengthMin:       helpers.GetIntP(d, "length_min"),

		HibpAllowedCount: helpers.GetIntP(d, "hibp_allowed_count"),

		ZxcvbnScoreThreshold: helpers.GetIntP(d, "zxcvbn_score_threshold"),
	}
	return &r
}

func resourcePolicyPasswordCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyPasswordSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordCreate(ctx).PasswordPolicyRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyPasswordRead(ctx, d, m)
}

func resourcePolicyPasswordRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "execution_logging", res.ExecutionLogging)
	helpers.SetWrapper(d, "password_field", res.PasswordField)
	helpers.SetWrapper(d, "error_message", res.ErrorMessage)
	helpers.SetWrapper(d, "amount_uppercase", res.AmountUppercase)
	helpers.SetWrapper(d, "amount_lowercase", res.AmountLowercase)
	helpers.SetWrapper(d, "amount_symbols", res.AmountSymbols)
	helpers.SetWrapper(d, "amount_digits", res.AmountDigits)
	helpers.SetWrapper(d, "length_min", res.LengthMin)
	helpers.SetWrapper(d, "symbol_charset", res.SymbolCharset)
	return diags
}

func resourcePolicyPasswordUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyPasswordSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesPasswordUpdate(ctx, d.Id()).PasswordPolicyRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyPasswordRead(ctx, d, m)
}

func resourcePolicyPasswordDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesPasswordDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
