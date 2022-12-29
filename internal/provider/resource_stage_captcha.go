package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceStageCaptcha() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageCaptchaCreate,
		ReadContext:   resourceStageCaptchaRead,
		UpdateContext: resourceStageCaptchaUpdate,
		DeleteContext: resourceStageCaptchaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"js_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceStageCaptchaSchemaToProvider(d *schema.ResourceData) *api.CaptchaStageRequest {
	r := api.CaptchaStageRequest{
		Name:       d.Get("name").(string),
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
	}
	if v, ok := d.GetOk("js_url"); ok {
		vv := v.(string)
		r.JsUrl = stringToPointer(vv)
	}
	if v, ok := d.GetOk("api_url"); ok {
		vv := v.(string)
		r.ApiUrl = stringToPointer(vv)
	}
	return &r
}

func resourceStageCaptchaCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageCaptchaSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesCaptchaCreate(ctx).CaptchaStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageCaptchaRead(ctx, d, m)
}

func resourceStageCaptchaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesCaptchaRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "public_key", res.PublicKey)
	setWrapper(d, "api_url", res.GetApiUrl())
	setWrapper(d, "js_url", res.GetJsUrl())
	return diags
}

func resourceStageCaptchaUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageCaptchaSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesCaptchaUpdate(ctx, d.Id()).CaptchaStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageCaptchaRead(ctx, d, m)
}

func resourceStageCaptchaDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesCaptchaDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
