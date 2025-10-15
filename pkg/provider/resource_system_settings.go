package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "days=365",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=30",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"default_token_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"reputation_lower_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -5,
			},
			"reputation_upper_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"flags": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          `{"policies_buffered_access_view": false}`,
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourceSystemSettingsSchemaToProvider(d *schema.ResourceData) (*api.SettingsRequest, diag.Diagnostics) {
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
		ReputationLowerLimit:      api.PtrInt32(int32(d.Get("reputation_lower_limit").(int))),
		ReputationUpperLimit:      api.PtrInt32(int32(d.Get("reputation_upper_limit").(int))),
	}

	flags, err := helpers.GetJSON[api.PatchedSettingsRequestFlags](d, ("flags"))
	r.Flags = flags
	return &r, err
}

func resourceSystemSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diag := resourceSystemSettingsSchemaToProvider(d)
	if diag != nil {
		return diag
	}

	_, hr, err := c.client.AdminApi.AdminSettingsUpdate(ctx).SettingsRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(systemSettingsID)
	return resourceSystemSettingsRead(ctx, d, m)
}

func resourceSystemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.AdminApi.AdminSettingsRetrieve(ctx).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "avatars", res.Avatars)
	helpers.SetWrapper(d, "default_user_change_name", res.DefaultUserChangeName)
	helpers.SetWrapper(d, "default_user_change_email", res.DefaultUserChangeEmail)
	helpers.SetWrapper(d, "default_user_change_username", res.DefaultUserChangeUsername)
	helpers.SetWrapper(d, "event_retention", res.EventRetention)
	helpers.SetWrapper(d, "footer_links", res.FooterLinks)
	helpers.SetWrapper(d, "gdpr_compliance", res.GdprCompliance)
	helpers.SetWrapper(d, "impersonation", res.Impersonation)
	helpers.SetWrapper(d, "default_token_duration", res.DefaultTokenDuration)
	helpers.SetWrapper(d, "default_token_length", res.DefaultTokenLength)
	helpers.SetWrapper(d, "reputation_lower_limit", res.ReputationLowerLimit)
	helpers.SetWrapper(d, "reputation_upper_limit", res.ReputationUpperLimit)
	return helpers.SetJSON(d, "flags", res.Flags)
}

func resourceSystemSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diag := resourceSystemSettingsSchemaToProvider(d)
	if diag != nil {
		return diag
	}

	_, hr, err := c.client.AdminApi.AdminSettingsUpdate(ctx).SettingsRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(systemSettingsID)
	return resourceSystemSettingsRead(ctx, d, m)
}
