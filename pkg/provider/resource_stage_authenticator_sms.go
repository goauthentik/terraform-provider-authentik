package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorSms() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sms_provider": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.PROVIDERENUM_TWILIO,
				Description:      helpers.EnumToDescription(api.AllowedProviderEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedProviderEnumEnumValues),
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
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHTYPEENUM_BASIC,
				Description:      helpers.EnumToDescription(api.AllowedAuthTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedAuthTypeEnumEnumValues),
			},
			"auth_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verify_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageAuthenticatorSmsSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorSMSStageRequest {
	r := api.AuthenticatorSMSStageRequest{
		Name:          d.Get("name").(string),
		Provider:      api.ProviderEnum(d.Get("sms_provider").(string)),
		FromNumber:    d.Get("from_number").(string),
		AccountSid:    d.Get("account_sid").(string),
		AuthType:      api.AuthTypeEnum(d.Get("auth_type").(string)).Ptr(),
		Auth:          d.Get("auth").(string),
		FriendlyName:  helpers.GetP[string](d, "friendly_name"),
		ConfigureFlow: *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),
		AuthPassword:  helpers.GetP[string](d, "auth_password"),
		Mapping:       *api.NewNullableString(helpers.GetP[string](d, "mapping")),
		VerifyOnly:    helpers.GetP[bool](d, "verify_only"),
	}
	return &r
}

func resourceStageAuthenticatorSmsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsCreate(ctx).AuthenticatorSMSStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "sms_provider", res.Provider)
	helpers.SetWrapper(d, "from_number", res.FromNumber)
	helpers.SetWrapper(d, "account_sid", res.AccountSid)
	helpers.SetWrapper(d, "auth", res.Auth)
	helpers.SetWrapper(d, "auth_password", res.AuthPassword)
	helpers.SetWrapper(d, "auth_type", res.AuthType)
	helpers.SetWrapper(d, "verify_only", res.VerifyOnly)
	helpers.SetWrapper(d, "mapping", res.Mapping.Get())
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	return diags
}

func resourceStageAuthenticatorSmsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorSmsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorSmsUpdate(ctx, d.Id()).AuthenticatorSMSStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorSmsRead(ctx, d, m)
}

func resourceStageAuthenticatorSmsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorSmsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
