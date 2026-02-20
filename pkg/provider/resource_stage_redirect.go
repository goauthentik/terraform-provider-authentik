package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageRedirect() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageRedirectCreate,
		ReadContext:   resourceStageRedirectRead,
		UpdateContext: resourceStageRedirectUpdate,
		DeleteContext: resourceStageRedirectDelete,
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
				Default:          api.REDIRECTSTAGEMODEENUM_FLOW,
				Description:      helpers.EnumToDescription(api.AllowedRedirectStageModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedRedirectStageModeEnumEnumValues),
			},
			"keep_context": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"target_static": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceStageRedirectSchemaToProvider(d *schema.ResourceData) *api.RedirectStageRequest {
	r := api.RedirectStageRequest{
		Name:         d.Get("name").(string),
		Mode:         api.RedirectStageModeEnum(d.Get("mode").(string)),
		KeepContext:  new(d.Get("keep_context").(bool)),
		TargetStatic: helpers.GetP[string](d, "target_static"),
		TargetFlow:   *api.NewNullableString(helpers.GetP[string](d, "target_flow")),
	}
	return &r
}

func resourceStageRedirectCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageRedirectSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesRedirectCreate(ctx).RedirectStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageRedirectRead(ctx, d, m)
}

func resourceStageRedirectRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesRedirectRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "mode", res.Mode)
	helpers.SetWrapper(d, "keep_context", res.KeepContext)
	helpers.SetWrapper(d, "target_flow", res.TargetFlow.Get())
	helpers.SetWrapper(d, "target_static", res.TargetStatic)
	return diags
}

func resourceStageRedirectUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageRedirectSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesRedirectUpdate(ctx, d.Id()).RedirectStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageRedirectRead(ctx, d, m)
}

func resourceStageRedirectDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesRedirectDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
