package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorEmail() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorEmailCreate,
		ReadContext:   resourceStageAuthenticatorEmailRead,
		UpdateContext: resourceStageAuthenticatorEmailUpdate,
		DeleteContext: resourceStageAuthenticatorEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_global_settings": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "localhost",
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"use_tls": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"use_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"from_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "system@authentik.local",
			},
			"token_expiry": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=30",
				Description:      RelativeDurationDescription,
				ValidateDiagFunc: ValidateRelativeDuration,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "authentik",
			},
			"template": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "email/password_reset.html",
			},
		},
	}
}

func resourceStageAuthenticatorEmailSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorEmailStageRequest {
	r := api.AuthenticatorEmailStageRequest{
		Name:              d.Get("name").(string),
		UseGlobalSettings: api.PtrBool(d.Get("use_global_settings").(bool)),
		UseSsl:            api.PtrBool(d.Get("use_ssl").(bool)),
		UseTls:            api.PtrBool(d.Get("use_tls").(bool)),
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
	}

	if h, hSet := d.GetOk("host"); hSet {
		r.Host = api.PtrString(h.(string))
	}
	if p, pSet := d.GetOk("port"); pSet {
		r.Port = api.PtrInt32(int32(p.(int)))
	}

	if h, hSet := d.GetOk("username"); hSet {
		r.Username = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("password"); hSet {
		r.Password = api.PtrString(h.(string))
	}

	if p, pSet := d.GetOk("timeout"); pSet {
		r.Timeout = api.PtrInt32(int32(p.(int)))
	}

	if h, hSet := d.GetOk("from_address"); hSet {
		r.FromAddress = api.PtrString(h.(string))
	}
	if p, pSet := d.GetOk("token_expiry"); pSet {
		r.TokenExpiry = api.PtrString(p.(string))
	}
	if h, hSet := d.GetOk("subject"); hSet {
		r.Subject = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("template"); hSet {
		r.Template = api.PtrString(h.(string))
	}
	return &r
}

func resourceStageAuthenticatorEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailCreate(ctx).AuthenticatorEmailStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEmailRead(ctx, d, m)
}

func resourceStageAuthenticatorEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "use_global_settings", res.UseGlobalSettings)
	setWrapper(d, "host", res.Host)
	setWrapper(d, "port", res.Port)
	setWrapper(d, "username", res.Username)
	setWrapper(d, "use_tls", res.UseTls)
	setWrapper(d, "use_ssl", res.UseSsl)
	setWrapper(d, "timeout", res.Timeout)
	setWrapper(d, "from_address", res.FromAddress)
	setWrapper(d, "token_expiry", res.TokenExpiry)
	setWrapper(d, "subject", res.Subject)
	setWrapper(d, "template", res.Template)
	setWrapper(d, "friendly_name", res.FriendlyName.Get())
	if res.ConfigureFlow.IsSet() {
		setWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	}
	return diags
}

func resourceStageAuthenticatorEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailUpdate(ctx, d.Id()).AuthenticatorEmailStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEmailRead(ctx, d, m)
}

func resourceStageAuthenticatorEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorEmailDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
