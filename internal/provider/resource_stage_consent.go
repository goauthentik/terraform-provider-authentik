package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStageConsent() *schema.Resource {
	return &schema.Resource{
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.CONSENTSTAGEMODEENUM_ALWAYS_REQUIRE,
			},
			"consent_expire_in": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "weeks=4",
			},
		},
	}
}

func resourceStageConsentSchemaToProvider(d *schema.ResourceData) (*api.ConsentStageRequest, diag.Diagnostics) {
	r := api.ConsentStageRequest{
		Name: d.Get("name").(string),
	}

	if m, mSet := d.GetOk("mode"); mSet {
		mo := api.ConsentStageModeEnum(m.(string))
		r.Mode = &mo
	}

	if ex, exSet := d.GetOk("consent_expire_in"); exSet {
		r.ConsentExpireIn = stringToPointer(ex.(string))
	}

	return &r, nil
}

func resourceStageConsentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageConsentSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesConsentCreate(ctx).ConsentStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesConsentRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("mode", res.Mode)
	d.Set("consent_expire_in", res.ConsentExpireIn)
	return diags
}

func resourceStageConsentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStageConsentSchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesConsentUpdate(ctx, d.Id()).ConsentStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageConsentRead(ctx, d, m)
}

func resourceStageConsentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesConsentDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
