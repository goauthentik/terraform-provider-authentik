package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
)

func resourceStageAuthenticatorWebAuthn() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageAuthenticatorWebAuthnCreate,
		ReadContext:   resourceStageAuthenticatorWebAuthnRead,
		UpdateContext: resourceStageAuthenticatorWebAuthnUpdate,
		DeleteContext: resourceStageAuthenticatorWebAuthnDelete,
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
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_verification": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERVERIFICATIONENUM_PREFERRED,
				Description:      helpers.EnumToDescription(api.AllowedUserVerificationEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserVerificationEnumEnumValues),
			},
			"resident_key_requirement": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.RESIDENTKEYREQUIREMENTENUM_PREFERRED,
				Description:      helpers.EnumToDescription(api.AllowedResidentKeyRequirementEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedResidentKeyRequirementEnumEnumValues),
			},
			"authenticator_attachment": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      helpers.EnumToDescription(api.AllowedAuthenticatorAttachmentEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedAuthenticatorAttachmentEnumEnumValues),
			},
			"device_type_restrictions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"max_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceStageAuthenticatorWebAuthnSchemaToProvider(d *schema.ResourceData) *api.AuthenticatorWebAuthnStageRequest {
	r := api.AuthenticatorWebAuthnStageRequest{
		Name:                   d.Get("name").(string),
		UserVerification:       api.UserVerificationEnum(d.Get("user_verification").(string)).Ptr(),
		ResidentKeyRequirement: api.ResidentKeyRequirementEnum(d.Get("resident_key_requirement").(string)).Ptr(),
		DeviceTypeRestrictions: castSlice[string](d.Get("device_type_restrictions").([]interface{})),
		FriendlyName:           *api.NewNullableString(getP[string](d, "friendly_name")),
		ConfigureFlow:          *api.NewNullableString(getP[string](d, "configure_flow")),
		MaxAttempts:            getIntP(d, "max_attempts"),
	}

	if x, set := d.GetOk("authenticator_attachment"); set {
		r.AuthenticatorAttachment.Set(api.AuthenticatorAttachmentEnum(x.(string)).Ptr())
	}
	return &r
}

func resourceStageAuthenticatorWebAuthnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnCreate(ctx).AuthenticatorWebAuthnStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "friendly_name", res.FriendlyName.Get())
	setWrapper(d, "user_verification", res.UserVerification)
	setWrapper(d, "resident_key_requirement", res.ResidentKeyRequirement)
	setWrapper(d, "authenticator_attachment", res.GetAuthenticatorAttachment())
	if res.ConfigureFlow.IsSet() {
		setWrapper(d, "configure_flow", res.ConfigureFlow.Get())
	}
	localDeviceTypeRestrictions := castSlice[string](d.Get("device_type_restrictions").([]interface{}))
	setWrapper(d, "device_type_restrictions", helpers.ListConsistentMerge(localDeviceTypeRestrictions, res.DeviceTypeRestrictions))
	setWrapper(d, "max_attempts", res.MaxAttempts)
	return diags
}

func resourceStageAuthenticatorWebAuthnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnUpdate(ctx, d.Id()).AuthenticatorWebAuthnStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
