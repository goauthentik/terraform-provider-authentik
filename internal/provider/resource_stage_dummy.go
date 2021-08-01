package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStageDummy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStageDummyCreate,
		ReadContext:   resourceStageDummyRead,
		UpdateContext: resourceStageDummyUpdate,
		DeleteContext: resourceStageDummyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceStageDummySchemaToProvider(d *schema.ResourceData) (*api.DummyStageRequest, diag.Diagnostics) {
	r := api.DummyStageRequest{
		Name: d.Get("name").(string),
	}
	return &r, nil
}

func resourceStageDummyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceStageDummySchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.StagesApi.StagesDummyCreate(ctx).DummyStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStageDummyRead(ctx, d, m)
}

func resourceStageDummyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesDummyRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	return diags
}

func resourceStageDummyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceStageDummySchemaToProvider(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.StagesApi.StagesDummyUpdate(ctx, d.Id()).DummyStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceStageDummyRead(ctx, d, m)
}

func resourceStageDummyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesDummyDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
