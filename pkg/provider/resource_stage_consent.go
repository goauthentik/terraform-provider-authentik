package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageConsent() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageConsentCreate,
		ReadContext:   resourceStageConsentRead,
		UpdateContext: resourceStageConsentUpdate,
		DeleteContext: resourceStageConsentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CONSENTSTAGEMODEENUM_ALWAYS_REQUIRE,
				Description:      helpers.EnumToDescription(api.AllowedConsentStageModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedConsentStageModeEnumEnumValues),
			},
			"consent_expire_in": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "weeks=4",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
		},
	}
}

func resourceStageConsentSchemaToProvider(d *schema.ResourceData) *api.ConsentStageRequest {
	r := api.ConsentStageRequest{
		Name:            d.Get("name").(string),
		ConsentExpireIn: helpers.GetP[string](d, "consent_expire_in"),
	}

	if m, mSet := d.GetOk("mode"); mSet {
		r.Mode = api.ConsentStageModeEnum(m.(string)).Ptr()
	}
	return &r
}

func resourceStageConsentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageConsentSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesConsentCreate(ctx).ConsentStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesConsentRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "mode", res.Mode)
	helpers.SetWrapper(d, "consent_expire_in", res.ConsentExpireIn)
	return diags
}

func resourceStageConsentUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageConsentSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesConsentUpdate(ctx, d.Id()).ConsentStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesConsentDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
