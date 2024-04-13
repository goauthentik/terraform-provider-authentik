package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      EnumToDescription(api.AllowedUserVerificationEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserVerificationEnumEnumValues),
			},
			"resident_key_requirement": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.RESIDENTKEYREQUIREMENTENUM_PREFERRED,
				Description:      EnumToDescription(api.AllowedResidentKeyRequirementEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedResidentKeyRequirementEnumEnumValues),
			},
			"authenticator_attachment": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      EnumToDescription(api.AllowedAuthenticatorAttachmentEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedAuthenticatorAttachmentEnumEnumValues),
			},
			"device_type_restrictions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
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
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
	setWrapper(d, "device_type_restrictions", listConsistentMerge(localDeviceTypeRestrictions, res.DeviceTypeRestrictions))
	return diags
}

func resourceStageAuthenticatorWebAuthnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnUpdate(ctx, d.Id()).AuthenticatorWebAuthnStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageAuthenticatorWebAuthnRead(ctx, d, m)
}

func resourceStageAuthenticatorWebAuthnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
