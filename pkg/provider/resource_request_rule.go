package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRequestRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceRequestRuleCreate,
		ReadContext:   resourceRequestRuleRead,
		UpdateContext: resourceRequestRuleUpdate,
		DeleteContext: resourceRequestRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"notification_transports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      helpers.EnumToDescription(api.AllowedNotificationModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedNotificationModeEnumEnumValues),
			},
			"min_reviewers": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_reviewers_is_per_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"request_flow": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID or slug of the flow shown to the requester when creating a request against this rule.",
			},
			// Computed
			"targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Objects currently bound to this rule, derived from its request rule bindings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRequestRuleSchemaToModel(d *schema.ResourceData) *api.RequestRuleRequest {
	m := api.RequestRuleRequest{
		Name:                   d.Get("name").(string),
		NotificationTransports: helpers.CastSlice[string](d, "notification_transports"),
		MinReviewers:           helpers.GetIntP(d, "min_reviewers"),
		MinReviewersIsPerGroup: helpers.GetP[bool](d, "min_reviewers_is_per_group"),
		RequestFlow:            *api.NewNullableString(helpers.GetP[string](d, "request_flow")),
	}
	if pem, ok := d.GetOk("policy_engine_mode"); ok {
		m.PolicyEngineMode = api.PolicyEngineMode(pem.(string)).Ptr()
	}
	if nm, ok := d.GetOk("notification_mode"); ok {
		m.NotificationMode = api.NotificationModeEnum(nm.(string)).Ptr()
	}
	return &m
}

func resourceRequestRuleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRulesCreate(ctx).RequestRuleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleRead(ctx, d, m)
}

func resourceRequestRuleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.RequestsAPI.RequestsRulesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "notification_transports", res.NotificationTransports)
	helpers.SetWrapper(d, "notification_mode", res.NotificationMode)
	helpers.SetWrapper(d, "min_reviewers", res.MinReviewers)
	helpers.SetWrapper(d, "min_reviewers_is_per_group", res.MinReviewersIsPerGroup)
	helpers.SetWrapper(d, "request_flow", res.RequestFlow.Get())
	helpers.SetWrapper(d, "targets", res.Targets)
	return diags
}

func resourceRequestRuleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceRequestRuleSchemaToModel(d)

	res, hr, err := c.client.RequestsAPI.RequestsRulesUpdate(ctx, d.Id()).RequestRuleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.GetUuid())
	return resourceRequestRuleRead(ctx, d, m)
}

func resourceRequestRuleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.RequestsAPI.RequestsRulesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
