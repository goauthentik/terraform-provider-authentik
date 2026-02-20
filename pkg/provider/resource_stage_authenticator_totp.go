package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorTOTP() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorTOTPCreate,
		ReadContext:   resourceStageAuthenticatorTOTPRead,
		UpdateContext: resourceStageAuthenticatorTOTPUpdate,
		DeleteContext: resourceStageAuthenticatorTOTPDelete,
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
			"digits": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.DIGITSENUM__6,
				Description:      helpers.EnumToDescription(api.AllowedDigitsEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedDigitsEnumEnumValues),
			},
		},
	}
}

func resourceStageAuthenticatorTOTPSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorTOTPStageRequest {
	r := api.AuthenticatorTOTPStageRequest{
		Name:          d.Get("name").(string),
		Digits:        api.DigitsEnum(d.Get("digits").(string)),
		FriendlyName:  helpers.GetP[string](d, "friendly_name"),
		ConfigureFlow: *api.NewNullableString(helpers.GetP[string](d, "configure_flow")),
	}
	return &r
}

func resourceStageAuthenticatorTOTPCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorTOTPSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpCreate(ctx).AuthenticatorTOTPStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorTOTPRead(ctx, d, m)
}

func resourceStageAuthenticatorTOTPRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "digits", res.Digits)
	helpers.SetWrapper(d, "friendly_name", res.FriendlyName)
	helpers.SetWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	return diags
}

func resourceStageAuthenticatorTOTPUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorTOTPSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorTotpUpdate(ctx, d.Id()).AuthenticatorTOTPStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorTOTPRead(ctx, d, m)
}

func resourceStageAuthenticatorTOTPDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorTotpDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
