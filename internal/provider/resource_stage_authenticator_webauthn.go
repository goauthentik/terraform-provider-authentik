package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageAuthenticatorWebAuthn() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"resident_key_requirement": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authenticator_attachment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageAuthenticatorWebAuthnSchemaToProvider(d *schema.ResourceData) *api.AuthenticateWebAuthnStageRequest {
	r := api.AuthenticateWebAuthnStageRequest{
		Name: d.Get("name").(string),
	}

	if fn, fnSet := d.GetOk("friendly_name"); fnSet {
		r.FriendlyName.Set(api.PtrString(fn.(string)))
	}
	if h, hSet := d.GetOk("configure_flow"); hSet {
		r.ConfigureFlow.Set(api.PtrString(h.(string)))
	}
	if x, set := d.GetOk("user_verification"); set {
		r.UserVerification = api.UserVerificationEnum(x.(string)).Ptr()
	}
	if x, set := d.GetOk("resident_key_requirement"); set {
		r.ResidentKeyRequirement = api.ResidentKeyRequirementEnum(x.(string)).Ptr()
	}
	if x, set := d.GetOk("authenticator_attachment"); set {
		r.AuthenticatorAttachment.Set(api.AuthenticatorAttachmentEnum(x.(string)).Ptr())
	}
	return &r
}

func resourceStageAuthenticatorWebAuthnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnCreate(ctx).AuthenticateWebAuthnStageRequest(*r).Execute()
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
	return diags
}

func resourceStageAuthenticatorWebAuthnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageAuthenticatorWebAuthnSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesAuthenticatorWebauthnUpdate(ctx, d.Id()).AuthenticateWebAuthnStageRequest(*app).Execute()
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
