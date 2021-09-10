package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowCreate,
		ReadContext:   resourceFlowRead,
		UpdateContext: resourceFlowUpdate,
		DeleteContext: resourceFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"designation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},
			"compatibility_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceFlowSchemaToModel(d *schema.ResourceData, c *APIClient) *api.FlowRequest {
	m := api.FlowRequest{
		Name:              d.Get("name").(string),
		Slug:              d.Get("slug").(string),
		Title:             d.Get("title").(string),
		Designation:       api.FlowDesignationEnum(d.Get("designation").(string)),
		CompatibilityMode: boolToPointer(d.Get("compatibility_mode").(bool)),
	}

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	m.PolicyEngineMode = &pm
	return &m
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d, c)

	res, hr, err := c.client.FlowsApi.FlowsInstancesCreate(ctx).FlowRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.FlowsApi.FlowsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("uuid", res.Pk)
	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	d.Set("title", res.Title)
	d.Set("designation", res.Designation)
	d.Set("policy_engine_mode", res.PolicyEngineMode)
	d.Set("compatibility_mode", res.CompatibilityMode)
	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d, c)

	res, hr, err := c.client.FlowsApi.FlowsInstancesUpdate(ctx, d.Id()).FlowRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.FlowsApi.FlowsInstancesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
