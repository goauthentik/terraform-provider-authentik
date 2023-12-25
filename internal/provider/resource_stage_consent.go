package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      EnumToDescription(api.AllowedConsentStageModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedConsentStageModeEnumEnumValues),
			},
			"consent_expire_in": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "weeks=4",
			},
		},
	}
}

func resourceStageConsentSchemaToProvider(d *schema.ResourceData) *api.ConsentStageRequest {
	r := api.ConsentStageRequest{
		Name: d.Get("name").(string),
	}

	if m, mSet := d.GetOk("mode"); mSet {
		r.Mode = api.ConsentStageModeEnum(m.(string)).Ptr()
	}

	if ex, exSet := d.GetOk("consent_expire_in"); exSet {
		r.ConsentExpireIn = api.PtrString(ex.(string))
	}
	return &r
}

func resourceStageConsentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageConsentSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesConsentCreate(ctx).ConsentStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesConsentRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "mode", res.Mode)
	setWrapper(d, "consent_expire_in", res.ConsentExpireIn)
	return diags
}

func resourceStageConsentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageConsentSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesConsentUpdate(ctx, d.Id()).ConsentStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesConsentDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
