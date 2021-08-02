package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFlowStageBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowStageBindingCreate,
		ReadContext:   resourceFlowStageBindingRead,
		UpdateContext: resourceFlowStageBindingUpdate,
		DeleteContext: resourceFlowStageBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluate_on_plan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"re_evaluate_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_engine_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.POLICYENGINEMODE_ANY,
			},
			"invalid_response_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.INVALIDRESPONSEACTIONENUM_RETRY,
			},
		},
	}
}

func resourceFlowStageBindingSchemaToModel(d *schema.ResourceData, c *APIClient) (*api.FlowStageBindingRequest, diag.Diagnostics) {
	m := api.FlowStageBindingRequest{
		Target:             d.Get("target").(string),
		Stage:              d.Get("stage").(string),
		Order:              int32(d.Get("order").(int)),
		EvaluateOnPlan:     boolToPointer(d.Get("evaluate_on_plan").(bool)),
		ReEvaluatePolicies: boolToPointer(d.Get("re_evaluate_policies").(bool)),
	}

	pm := api.PolicyEngineMode(d.Get("policy_engine_mode").(string))
	m.PolicyEngineMode = &pm

	ira := api.InvalidResponseActionEnum(d.Get("invalid_response_action").(string))
	m.InvalidResponseAction = &ira
	return &m, nil
}

func resourceFlowStageBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceFlowStageBindingSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.FlowsApi.FlowsBindingsCreate(ctx).FlowStageBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceFlowStageBindingRead(ctx, d, m)
}

func resourceFlowStageBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.FlowsApi.FlowsBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("target", res.Target)
	d.Set("stage", res.Stage)
	d.Set("order", res.Order)
	d.Set("evaluate_on_plan", res.EvaluateOnPlan)
	d.Set("re_evaluate_policies", res.ReEvaluatePolicies)
	d.Set("policy_engine_mode", res.PolicyEngineMode)
	d.Set("invalid_response_action", res.InvalidResponseAction)
	return diags
}

func resourceFlowStageBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceFlowStageBindingSchemaToModel(d, c)
	if di != nil {
		return di
	}

	res, hr, err := c.client.FlowsApi.FlowsBindingsUpdate(ctx, d.Id()).FlowStageBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Pk)
	return resourceFlowStageBindingRead(ctx, d, m)
}

func resourceFlowStageBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.FlowsApi.FlowsBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
