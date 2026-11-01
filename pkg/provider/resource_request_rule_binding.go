package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRequestRuleBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceRequestRuleBindingCreate,
		ReadContext:   resourceRequestRuleBindingRead,
		UpdateContext: resourceRequestRuleBindingUpdate,
		DeleteContext: resourceRequestRuleBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the authentik_request_rule this binding applies.",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the object this binding makes requestable.",
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"expiry_pending": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "How long a request against this binding stays pending before it automatically lapses if not approved or denied. Format: hours=1;minutes=2;seconds=3.",
			},
			"expiry_granted_max": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The maximum duration a grant approved against this binding can last. Format: hours=1;minutes=2;seconds=3.",
			},
			// Computed
			"related": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Objects related to this binding's target, derived server-side.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRequestRuleBindingSchemaToModel(d *schema.ResourceData) *api.RequestRuleBindingRequest {
	m := api.RequestRuleBindingRequest{
		Rule:             d.Get("rule").(string),
		Target:           d.Get("target").(string),
		ExpiryPending:    helpers.GetP[string](d, "expiry_pending"),
		ExpiryGrantedMax: helpers.GetP[string](d, "expiry_granted_max"),
	}
	if pem, ok := d.GetOk("policy_engine_mode"); ok {
		m.PolicyEngineMode = api.PolicyEngineMode(pem.(string)).Ptr()
	}
	return &m
}

func resourceRequestRuleBindingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleBindingSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRuleBindingsCreate(ctx).RequestRuleBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleBindingRead(ctx, d, m)
}

func resourceRequestRuleBindingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RequestsAPI.RequestsRuleBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "rule", res.Rule)
	helpers.SetWrapper(d, "target", res.Target)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "expiry_pending", res.ExpiryPending)
	helpers.SetWrapper(d, "expiry_granted_max", res.ExpiryGrantedMax)
	helpers.SetWrapper(d, "related", res.Related)
	return diags
}

func resourceRequestRuleBindingUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleBindingSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRuleBindingsUpdate(ctx, d.Id()).RequestRuleBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleBindingRead(ctx, d, m)
}

func resourceRequestRuleBindingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RequestsAPI.RequestsRuleBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
