package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceEventRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Events --- ",
		CreateContext: resourceEventRuleCreate,
		ReadContext:   resourceEventRuleRead,
		UpdateContext: resourceEventRuleUpdate,
		DeleteContext: resourceEventRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transports": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"severity": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.SEVERITYENUM_WARNING,
				Description:      helpers.EnumToDescription(api.AllowedSeverityEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedSeverityEnumEnumValues),
			},
			"destination_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group to send notification to",
			},
			"destination_event_user": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Send notification to event user",
			},
		},
	}
}

func resourceEventRuleSchemaToModel(d *schema.ResourceData) (*api.NotificationRuleRequest, diag.Diagnostics) {
	m := api.NotificationRuleRequest{
		Name:                 d.Get("name").(string),
		Severity:             api.SeverityEnum(d.Get("severity").(string)).Ptr(),
		DestinationEventUser: new(d.Get("destination_event_user").(bool)),
		DestinationGroup:     *api.NewNullableString(helpers.GetP[string](d, ("destination_group"))),
		Transports:           helpers.CastSlice[string](d, "transports"),
	}
	return &m, nil
}

func resourceEventRuleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEventRuleSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EventsApi.EventsRulesCreate(ctx).NotificationRuleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventRuleRead(ctx, d, m)
}

func resourceEventRuleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EventsApi.EventsRulesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "destination_group", res.DestinationGroup.Get())
	helpers.SetWrapper(d, "destination_event_user", res.DestinationEventUser)
	helpers.SetWrapper(d, "transports", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "transports"),
		res.Transports,
	))
	helpers.SetWrapper(d, "severity", res.Severity)
	return diags
}

func resourceEventRuleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceEventRuleSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.EventsApi.EventsRulesUpdate(ctx, d.Id()).NotificationRuleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceEventRuleRead(ctx, d, m)
}

func resourceEventRuleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EventsApi.EventsRulesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
