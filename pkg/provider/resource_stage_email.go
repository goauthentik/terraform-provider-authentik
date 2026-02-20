package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
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
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
		},
	}
}

func resourceStageEmailSchemaToProvider(d *schema.ResourceData) *api.EmailStageRequest {
	r := api.EmailStageRequest{
		Name:                  d.Get("name").(string),
		UseGlobalSettings:     new(d.Get("use_global_settings").(bool)),
		UseSsl:                new(d.Get("use_ssl").(bool)),
		UseTls:                new(d.Get("use_tls").(bool)),
		Host:                  helpers.GetP[string](d, "host"),
		Username:              helpers.GetP[string](d, "username"),
		Password:              helpers.GetP[string](d, "password"),
		FromAddress:           helpers.GetP[string](d, "from_address"),
		TokenExpiry:           helpers.GetP[string](d, "token_expiry"),
		Subject:               helpers.GetP[string](d, "subject"),
		Template:              helpers.GetP[string](d, "template"),
		RecoveryCacheTimeout:  helpers.GetP[string](d, "recovery_cache_timeout"),
		Port:                  helpers.GetIntP(d, "port"),
		Timeout:               helpers.GetIntP(d, "timeout"),
		RecoveryMaxAttempts:   helpers.GetIntP(d, "recovery_max_attempts"),
		ActivateUserOnSuccess: helpers.GetP[bool](d, "activate_user_on_success"),
	}
	return &r
}

func resourceStageEmailCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEmailCreate(ctx).EmailStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEmailRead(ctx, d, m)
}

func resourceStageEmailRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesEmailRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "use_global_settings", res.UseGlobalSettings)
	helpers.SetWrapper(d, "host", res.Host)
	helpers.SetWrapper(d, "port", res.Port)
	helpers.SetWrapper(d, "username", res.Username)
	helpers.SetWrapper(d, "use_tls", res.UseTls)
	helpers.SetWrapper(d, "use_ssl", res.UseSsl)
	helpers.SetWrapper(d, "timeout", res.Timeout)
	helpers.SetWrapper(d, "from_address", res.FromAddress)
	helpers.SetWrapper(d, "token_expiry", res.TokenExpiry)
	helpers.SetWrapper(d, "subject", res.Subject)
	helpers.SetWrapper(d, "template", res.Template)
	helpers.SetWrapper(d, "activate_user_on_success", res.ActivateUserOnSuccess)
	helpers.SetWrapper(d, "recovery_max_attempts", res.RecoveryMaxAttempts)
	helpers.SetWrapper(d, "recovery_cache_timeout", res.RecoveryCacheTimeout)
	return diags
}

func resourceStageEmailUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEmailUpdate(ctx, d.Id()).EmailStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEmailRead(ctx, d, m)
}

func resourceStageEmailDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesEmailDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
