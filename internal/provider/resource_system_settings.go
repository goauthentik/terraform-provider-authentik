package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

const systemSettingsID = "system_settings"

func resourceSystemSettings() *schema.Resource {
	return &schema.Resource{
		Description:   "System --- ",
		CreateContext: resourceSystemSettingsCreate,
		ReadContext:   resourceSystemSettingsRead,
		UpdateContext: resourceSystemSettingsUpdate,
		DeleteContext: schema.NoopContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"avatars": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "gravatar,initials",
			},
			"default_user_change_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"default_user_change_email": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"default_user_change_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"event_retention": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=365",
			},
			"footer_links": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"gdpr_compliance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"impersonation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"default_token_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "minutes=30",
			},
			"default_token_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func resourceSystemSettingsSchemaToProvider(d *schema.ResourceData) *api.SettingsRequest {
	r := api.SettingsRequest{
		Avatars:                   api.PtrString(d.Get("avatars").(string)),
		DefaultUserChangeName:     api.PtrBool(d.Get("default_user_change_name").(bool)),
		DefaultUserChangeEmail:    api.PtrBool(d.Get("default_user_change_email").(bool)),
		DefaultUserChangeUsername: api.PtrBool(d.Get("default_user_change_username").(bool)),
		EventRetention:            api.PtrString(d.Get("event_retention").(string)),
		FooterLinks:               d.Get("footer_links"),
		GdprCompliance:            api.PtrBool(d.Get("gdpr_compliance").(bool)),
		Impersonation:             api.PtrBool(d.Get("impersonation").(bool)),
		DefaultTokenDuration:      api.PtrString(d.Get("default_token_duration").(string)),
		DefaultTokenLength:        api.PtrInt32(int32(d.Get("default_token_length").(int))),
	}
	return &r
}

func resourceSystemSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSystemSettingsSchemaToProvider(d)

	_, hr, err := c.client.AdminApi.AdminSettingsUpdate(ctx).SettingsRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(systemSettingsID)
	return resourceSystemSettingsRead(ctx, d, m)
}

func resourceSystemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.AdminApi.AdminSettingsRetrieve(ctx).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "avatars", res.Avatars)
	setWrapper(d, "default_user_change_name", res.DefaultUserChangeName)
	setWrapper(d, "default_user_change_email", res.DefaultUserChangeEmail)
	setWrapper(d, "default_user_change_username", res.DefaultUserChangeUsername)
	setWrapper(d, "event_retention", res.EventRetention)
	setWrapper(d, "footer_links", res.FooterLinks)
	setWrapper(d, "gdpr_compliance", res.GdprCompliance)
	setWrapper(d, "impersonation", res.Impersonation)
	setWrapper(d, "default_token_duration", res.DefaultTokenDuration)
	setWrapper(d, "default_token_length", res.DefaultTokenLength)
	return diags
}

func resourceSystemSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceSystemSettingsSchemaToProvider(d)

	_, hr, err := c.client.AdminApi.AdminSettingsUpdate(ctx).SettingsRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(systemSettingsID)
	return resourceSystemSettingsRead(ctx, d, m)
}
