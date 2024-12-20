package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceFlowStageBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Evaluate policies during the Flow planning process.",
			},
			"re_evaluate_policies": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Evaluate policies when the Stage is present to the user.",
			},
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"invalid_response_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.INVALIDRESPONSEACTIONENUM_RETRY,
				Description:      EnumToDescription(api.AllowedInvalidResponseActionEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedInvalidResponseActionEnumEnumValues),
			},
		},
	}
}

func resourceFlowStageBindingSchemaToModel(d *schema.ResourceData) *api.FlowStageBindingRequest {
	m := api.FlowStageBindingRequest{
		Target:                d.Get("target").(string),
		Stage:                 d.Get("stage").(string),
		Order:                 int32(d.Get("order").(int)),
		EvaluateOnPlan:        api.PtrBool(d.Get("evaluate_on_plan").(bool)),
		ReEvaluatePolicies:    api.PtrBool(d.Get("re_evaluate_policies").(bool)),
		PolicyEngineMode:      api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		InvalidResponseAction: api.InvalidResponseActionEnum(d.Get("invalid_response_action").(string)).Ptr(),
	}
	return &m
}

func resourceFlowStageBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowStageBindingSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsBindingsCreate(ctx).FlowStageBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceFlowStageBindingRead(ctx, d, m)
}

func resourceFlowStageBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.FlowsApi.FlowsBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "target", res.Target)
	setWrapper(d, "stage", res.Stage)
	setWrapper(d, "order", res.Order)
	setWrapper(d, "evaluate_on_plan", res.EvaluateOnPlan)
	setWrapper(d, "re_evaluate_policies", res.ReEvaluatePolicies)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "invalid_response_action", res.InvalidResponseAction)
	return diags
}

func resourceFlowStageBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowStageBindingSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsBindingsUpdate(ctx, d.Id()).FlowStageBindingRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceFlowStageBindingRead(ctx, d, m)
}

func resourceFlowStageBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.FlowsApi.FlowsBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
