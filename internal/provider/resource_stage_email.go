package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageEmail() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageEmailCreate,
		ReadContext:   resourceStageEmailRead,
		UpdateContext: resourceStageEmailUpdate,
		DeleteContext: resourceStageEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
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
			"activate_user_on_success": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageEmailSchemaToProvider(d *schema.ResourceData) *api.EmailStageRequest {
	r := api.EmailStageRequest{
		Name:              d.Get("name").(string),
		UseGlobalSettings: api.PtrBool(d.Get("use_global_settings").(bool)),
		UseSsl:            api.PtrBool(d.Get("use_ssl").(bool)),
		UseTls:            api.PtrBool(d.Get("use_tls").(bool)),
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
		r.TokenExpiry = api.PtrInt32(int32(p.(int)))
	}
	if h, hSet := d.GetOk("subject"); hSet {
		r.Subject = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("template"); hSet {
		r.Template = api.PtrString(h.(string))
	}
	if h, hSet := d.GetOk("activate_user_on_success"); hSet {
		r.ActivateUserOnSuccess = api.PtrBool(h.(bool))
	}
	return &r
}

func resourceStageEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEmailCreate(ctx).EmailStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEmailRead(ctx, d, m)
}

func resourceStageEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesEmailRetrieve(ctx, d.Id()).Execute()
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
	setWrapper(d, "activate_user_on_success", res.ActivateUserOnSuccess)
	return diags
}

func resourceStageEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEmailUpdate(ctx, d.Id()).EmailStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEmailRead(ctx, d, m)
}

func resourceStageEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesEmailDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
