package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageAuthenticatorValidate() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorValidateCreate,
		ReadContext:   resourceStageAuthenticatorValidateRead,
		UpdateContext: resourceStageAuthenticatorValidateUpdate,
		DeleteContext: resourceStageAuthenticatorValidateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"not_configured_action": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedNotConfiguredActionEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedNotConfiguredActionEnumEnumValues),
			},
			"device_classes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					Description:      helpers.EnumToDescription(api.AllowedDeviceClassesEnumEnumValues),
					ValidateDiagFunc: helpers.StringInEnum(api.AllowedDeviceClassesEnumEnumValues),
				},
			},
			"configuration_stages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_auth_threshold": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "seconds=0",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
			"webauthn_user_verification": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERVERIFICATIONENUM_PREFERRED,
				Description:      helpers.EnumToDescription(api.AllowedUserVerificationEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserVerificationEnumEnumValues),
			},
			"webauthn_allowed_device_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStageAuthenticatorValidateSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorValidateStageRequest {
	r := api.AuthenticatorValidateStageRequest{
		Name:                       d.Get("name").(string),
		LastAuthThreshold:          new(d.Get("last_auth_threshold").(string)),
		WebauthnAllowedDeviceTypes: helpers.CastSlice[string](d, "webauthn_allowed_device_types"),
		NotConfiguredAction:        helpers.GetP[api.NotConfiguredActionEnum](d, "not_configured_action"),
		ConfigurationStages:        helpers.CastSlice[string](d, "configuration_stages"),
		WebauthnUserVerification:   helpers.CastString[api.UserVerificationEnum](helpers.GetP[string](d, "webauthn_user_verification")),
	}

	classes := make([]api.DeviceClassesEnum, 0)
	for _, classesS := range d.Get("device_classes").([]any) {
		classes = append(classes, api.DeviceClassesEnum(classesS.(string)))
	}
	r.DeviceClasses = classes
	return &r
}

func resourceStageAuthenticatorValidateCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateCreate(ctx).AuthenticatorValidateStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "not_configured_action", res.NotConfiguredAction)
	helpers.SetWrapper(d, "configuration_stages", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "configuration_stages"),
		res.ConfigurationStages,
	))
	helpers.SetWrapper(d, "device_classes", res.DeviceClasses)
	helpers.SetWrapper(d, "last_auth_threshold", res.LastAuthThreshold)
	helpers.SetWrapper(d, "webauthn_user_verification", res.WebauthnUserVerification)
	helpers.SetWrapper(d, "webauthn_allowed_device_types", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "webauthn_allowed_device_types"),
		res.WebauthnAllowedDeviceTypes,
	))
	return diags
}

func resourceStageAuthenticatorValidateUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateUpdate(ctx, d.Id()).AuthenticatorValidateStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorValidateDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
