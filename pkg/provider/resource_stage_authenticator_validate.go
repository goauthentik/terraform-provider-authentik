package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      EnumToDescription(api.AllowedNotConfiguredActionEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedNotConfiguredActionEnumEnumValues),
			},
			"device_classes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					Description:      EnumToDescription(api.AllowedDeviceClassesEnumEnumValues),
					ValidateDiagFunc: StringInEnum(api.AllowedDeviceClassesEnumEnumValues),
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  "seconds=0",
			},
			"webauthn_user_verification": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERVERIFICATIONENUM_PREFERRED,
				Description:      EnumToDescription(api.AllowedUserVerificationEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserVerificationEnumEnumValues),
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
		Name:              d.Get("name").(string),
		LastAuthThreshold: api.PtrString(d.Get("last_auth_threshold").(string)),
	}

	if h, hSet := d.GetOk("not_configured_action"); hSet {
		r.NotConfiguredAction = api.NotConfiguredActionEnum(h.(string)).Ptr()
	}
	if h, hSet := d.GetOk("configuration_stages"); hSet {
		r.ConfigurationStages = castSlice[string](h.([]interface{}))
	}
	if x, set := d.GetOk("webauthn_user_verification"); set {
		r.WebauthnUserVerification = api.UserVerificationEnum(x.(string)).Ptr()
	}

	classes := make([]api.DeviceClassesEnum, 0)
	for _, classesS := range d.Get("device_classes").([]interface{}) {
		classes = append(classes, api.DeviceClassesEnum(classesS.(string)))
	}
	r.DeviceClasses = classes
	r.WebauthnAllowedDeviceTypes = castSlice[string](d.Get("webauthn_allowed_device_types").([]interface{}))
	return &r
}

func resourceStageAuthenticatorValidateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateCreate(ctx).AuthenticatorValidateStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "not_configured_action", res.NotConfiguredAction)
	if res.ConfigurationStages != nil {
		localConfigurationStages := castSlice[string](d.Get("configuration_stages").([]interface{}))
		setWrapper(d, "configuration_stages", listConsistentMerge(localConfigurationStages, res.ConfigurationStages))
	}
	setWrapper(d, "device_classes", res.DeviceClasses)
	setWrapper(d, "last_auth_threshold", res.LastAuthThreshold)
	setWrapper(d, "webauthn_user_verification", res.WebauthnUserVerification)
	localDeviceTypeRestrictions := castSlice[string](d.Get("webauthn_allowed_device_types").([]interface{}))
	setWrapper(d, "webauthn_allowed_device_types", listConsistentMerge(localDeviceTypeRestrictions, res.WebauthnAllowedDeviceTypes))
	return diags
}

func resourceStageAuthenticatorValidateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorValidateSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorValidateUpdate(ctx, d.Id()).AuthenticatorValidateStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorValidateRead(ctx, d, m)
}

func resourceStageAuthenticatorValidateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorValidateDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
