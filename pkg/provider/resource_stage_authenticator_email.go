package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Default:  "",
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
		},
	}
}

func resourceStageAuthenticatorEmailSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorEmailStageRequest {
	r := api.AuthenticatorEmailStageRequest{
		Name:              d.Get("name").(string),
		UseGlobalSettings: new(d.Get("use_global_settings").(bool)),
		UseSsl:            new(d.Get("use_ssl").(bool)),
		UseTls:            new(d.Get("use_tls").(bool)),
		FriendlyName:      helpers.GetP[string](d, "friendly_name"),
		ConfigureFlow:     *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),

		Host: helpers.GetP[string](d, "host"),
		Port: helpers.GetIntP(d, "port"),

		Username: helpers.GetP[string](d, "username"),
		Password: helpers.GetP[string](d, "password"),

		Timeout: helpers.GetIntP(d, "timeout"),

		FromAddress: helpers.GetP[string](d, "from_address"),
		TokenExpiry: helpers.GetP[string](d, "token_expiry"),
		Subject:     helpers.GetP[string](d, "subject"),
		Template:    helpers.GetP[string](d, "template"),
	}
	return &r
}

func resourceStageAuthenticatorEmailCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailCreate(ctx).AuthenticatorEmailStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEmailRead(ctx, d, m)
}

func resourceStageAuthenticatorEmailRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailRetrieve(ctx, d.Id()).Execute()
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
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	return diags
}

func resourceStageAuthenticatorEmailUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorEmailSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorEmailUpdate(ctx, d.Id()).AuthenticatorEmailStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorEmailRead(ctx, d, m)
}

func resourceStageAuthenticatorEmailDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorEmailDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
