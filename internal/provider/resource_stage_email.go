package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceStageEmail() *schema.Resource {
	return &schema.Resource{
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
		},
	}
}

func resourceStageEmailSchemaToProvider(d *schema.ResourceData) *api.EmailStageRequest {
	r := api.EmailStageRequest{
		Name:              d.Get("name").(string),
		UseGlobalSettings: boolToPointer(d.Get("use_global_settings").(bool)),
		UseSsl:            boolToPointer(d.Get("use_ssl").(bool)),
		UseTls:            boolToPointer(d.Get("use_tls").(bool)),
	}

	if h, hSet := d.GetOk("host"); hSet {
		r.Host = stringToPointer(h.(string))
	}
	if p, pSet := d.GetOk("port"); pSet {
		r.Port = intToPointer(p.(int))
	}

	if h, hSet := d.GetOk("username"); hSet {
		r.Username = stringToPointer(h.(string))
	}
	if h, hSet := d.GetOk("password"); hSet {
		r.Password = stringToPointer(h.(string))
	}

	if p, pSet := d.GetOk("timeout"); pSet {
		r.Timeout = intToPointer(p.(int))
	}

	if h, hSet := d.GetOk("from_address"); hSet {
		r.FromAddress = stringToPointer(h.(string))
	}
	if p, pSet := d.GetOk("token_expiry"); pSet {
		r.TokenExpiry = intToPointer(p.(int))
	}
	if h, hSet := d.GetOk("subject"); hSet {
		r.Subject = stringToPointer(h.(string))
	}
	if h, hSet := d.GetOk("template"); hSet {
		r.Template = stringToPointer(h.(string))
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

	d.Set("name", res.Name)
	d.Set("use_global_settings", res.UseGlobalSettings)
	d.Set("host", res.Host)
	d.Set("port", res.Port)
	d.Set("username", res.Username)
	d.Set("use_tls", res.UseTls)
	d.Set("use_ssl", res.UseSsl)
	d.Set("timeout", res.Timeout)
	d.Set("from_address", res.FromAddress)
	d.Set("token_expiry", res.TokenExpiry)
	d.Set("subject", res.Subject)
	d.Set("template", res.Template)
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
