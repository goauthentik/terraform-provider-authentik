package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageCaptcha() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
				Default:  "https://www.recaptcha.net/recaptcha/api.js",
			},
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://www.recaptcha.net/recaptcha/api/siteverify",
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"score_min_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  1,
			},
			"score_max_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0.5,
			},
			"error_on_invalid_score": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"interactive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceStageCaptchaSchemaToProvider(d *schema.ResourceData) *api.CaptchaStageRequest {
	r := api.CaptchaStageRequest{
		Name:                d.Get("name").(string),
		PublicKey:           d.Get("public_key").(string),
		PrivateKey:          d.Get("private_key").(string),
		ErrorOnInvalidScore: new(d.Get("error_on_invalid_score").(bool)),
		ScoreMinThreshold:   new(d.Get("score_min_threshold").(float64)),
		ScoreMaxThreshold:   new(d.Get("score_max_threshold").(float64)),
		Interactive:         new(d.Get("interactive").(bool)),
		JsUrl:               helpers.GetP[string](d, "js_url"),
		ApiUrl:              helpers.GetP[string](d, "api_url"),
	}
	return &r
}

func resourceStageCaptchaCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageCaptchaSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesCaptchaCreate(ctx).CaptchaStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageCaptchaRead(ctx, d, m)
}

func resourceStageCaptchaRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesCaptchaRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "public_key", res.PublicKey)
	helpers.SetWrapper(d, "api_url", res.GetApiUrl())
	helpers.SetWrapper(d, "js_url", res.GetJsUrl())
	helpers.SetWrapper(d, "error_on_invalid_score", res.GetErrorOnInvalidScore())
	helpers.SetWrapper(d, "score_min_threshold", res.GetScoreMinThreshold())
	helpers.SetWrapper(d, "score_max_threshold", res.GetScoreMaxThreshold())
	helpers.SetWrapper(d, "interactive", res.Interactive)
	return diags
}

func resourceStageCaptchaUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageCaptchaSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesCaptchaUpdate(ctx, d.Id()).CaptchaStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageCaptchaRead(ctx, d, m)
}

func resourceStageCaptchaDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesCaptchaDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
