package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceStageAuthenticatorSms() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageAuthenticatorSmsCreate,
		ReadContext:   resourceStageAuthenticatorSmsRead,
		UpdateContext: resourceStageAuthenticatorSmsUpdate,
		DeleteContext: resourceStageAuthenticatorSmsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sms_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.PROVIDERENUM_TWILIO,
			},
			"from_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_sid": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.AUTHTYPEENUM_BASIC,
			},
			"auth_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceStageAuthenticatorSmsSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorSMSStageRequest {
	r := api.AuthenticatorSMSStageRequest{
		Name:       d.Get("name").(string),
		Provider:   api.ProviderEnum(d.Get("sms_provider").(string)),
		FromNumber: d.Get("from_number").(string),
		AccountSid: d.Get("account_sid").(string),
		Auth:       d.Get("auth").(string),
	}

	at := api.AuthTypeEnum(d.Get("auth_type").(string))
	r.AuthType = &at

	if h, hSet := d.GetOk("auth_password"); hSet {
		r.AuthPassword = stringToPointer(h.(string))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(stringToPointer(h.(string)))
	}
	return &r
}

func resourceStageAuthenticatorSmsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsCreate(ctx).AuthenticatorSMSStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("sms_provider", res.Provider)
	d.Set("from_number", res.FromNumber)
	d.Set("account_sid", res.AccountSid)
	d.Set("auth", res.Auth)
	d.Set("auth_password", res.AuthPassword)
	d.Set("auth_type", res.AuthType)
	if res.ConfigureFlow.IsSet() {
		d.Set("configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorSmsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsUpdate(ctx, d.Id()).AuthenticatorSMSStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorSmsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
