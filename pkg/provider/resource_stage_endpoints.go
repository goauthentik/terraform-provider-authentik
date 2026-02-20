package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageEndpoints() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageEndpointsCreate,
		ReadContext:   resourceStageEndpointsRead,
		UpdateContext: resourceStageEndpointsUpdate,
		DeleteContext: resourceStageEndpointsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connector": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.STAGEMODEENUM_OPTIONAL,
				Description:      helpers.EnumToDescription(api.AllowedStageModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedStageModeEnumEnumValues),
			},
		},
	}
}

func resourceStageEndpointsSchemaToProvider(d *schema.ResourceData) *api.EndpointStageRequest {
	r := api.EndpointStageRequest{
		Name:      d.Get("name").(string),
		Connector: d.Get("connector").(string),
		Mode:      api.StageModeEnum(d.Get("mode").(string)).Ptr(),
	}
	return &r
}

func resourceStageEndpointsCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageEndpointsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEndpointsCreate(ctx).EndpointStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEndpointsRead(ctx, d, m)
}

func resourceStageEndpointsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesEndpointsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "connector", res.Connector)
	helpers.SetWrapper(d, "mode", res.Mode)
	return diags
}

func resourceStageEndpointsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageEndpointsSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesEndpointsUpdate(ctx, d.Id()).EndpointStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageEndpointsRead(ctx, d, m)
}

func resourceStageEndpointsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesEndpointsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
