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
			"activate_user_on_success": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"recovery_max_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"recovery_cache_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=5",
				Description:      RelativeDurationDescription,
				ValidateDiagFunc: ValidateRelativeDuration,
			},
		},
	}
}

func resourceStageEmailSchemaToProvider(d *schema.ResourceData) *api.EmailStageRequest {
	r := api.EmailStageRequest{
		Name:                  d.Get("name").(string),
		UseGlobalSettings:     api.PtrBool(d.Get("use_global_settings").(bool)),
		UseSsl:                api.PtrBool(d.Get("use_ssl").(bool)),
		UseTls:                api.PtrBool(d.Get("use_tls").(bool)),
		Host:                  getP[string](d, "host"),
		Username:              getP[string](d, "username"),
		Password:              getP[string](d, "password"),
		FromAddress:           getP[string](d, "from_address"),
		TokenExpiry:           getP[string](d, "token_expiry"),
		Subject:               getP[string](d, "subject"),
		Template:              getP[string](d, "template"),
		RecoveryCacheTimeout:  getP[string](d, "recovery_cache_timeout"),
		Port:                  getIntP(d, "port"),
		Timeout:               getIntP(d, "timeout"),
		RecoveryMaxAttempts:   getIntP(d, "recovery_max_attempts"),
		ActivateUserOnSuccess: getP[bool](d, "activate_user_on_success"),
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
	setWrapper(d, "recovery_max_attempts", res.RecoveryMaxAttempts)
	setWrapper(d, "recovery_cache_timeout", res.RecoveryCacheTimeout)
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
