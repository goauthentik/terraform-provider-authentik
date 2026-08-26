package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// resourceUserOffboarding has no Update endpoint in the API - every field is ForceNew and
// UpdateContext is intentionally omitted, so any config change recreates the schedule.
func resourceUserOffboarding() *schema.Resource {
	return &schema.Resource{
		Description:   "Enterprise --- ",
		CreateContext: resourceUserOffboardingCreate,
		ReadContext:   resourceUserOffboardingRead,
		DeleteContext: resourceUserOffboardingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "PK of the user to offboard.",
			},
			"scheduled_at": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Absolute RFC3339 time at which the offboarding action is executed.",
			},
			"action": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Description:      helpers.EnumToDescription(api.AllowedOffboardingActionEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedOffboardingActionEnumEnumValues),
			},
			"revoke_sessions": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"revoke_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"executed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUserOffboardingSchemaToModel(d *schema.ResourceData) (*api.UserOffboardingRequest, diag.Diagnostics) {
	scheduledAt, err := time.Parse(time.RFC3339, d.Get("scheduled_at").(string))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	m := api.UserOffboardingRequest{
		User:           int32(d.Get("user").(int)),
		ScheduledAt:    scheduledAt,
		RevokeSessions: helpers.GetP[bool](d, "revoke_sessions"),
		RevokeTokens:   helpers.GetP[bool](d, "revoke_tokens"),
	}
	if a, ok := d.GetOk("action"); ok {
		m.Action = api.OffboardingActionEnum(a.(string)).Ptr()
	}
	return &m, nil
}

func resourceUserOffboardingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceUserOffboardingSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.LifecycleAPI.LifecycleUserOffboardingCreate(ctx).UserOffboardingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Id)
	return resourceUserOffboardingRead(ctx, d, m)
}

func resourceUserOffboardingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.LifecycleAPI.LifecycleUserOffboardingRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "user", res.User)
	helpers.SetWrapper(d, "scheduled_at", res.ScheduledAt.Format(time.RFC3339))
	helpers.SetWrapper(d, "action", res.Action)
	helpers.SetWrapper(d, "revoke_sessions", res.RevokeSessions)
	helpers.SetWrapper(d, "revoke_tokens", res.RevokeTokens)
	helpers.SetWrapper(d, "status", res.Status)
	helpers.SetWrapper(d, "created_at", res.CreatedAt.Format(time.RFC3339))
	if res.ExecutedAt.IsSet() && res.ExecutedAt.Get() != nil {
		helpers.SetWrapper(d, "executed_at", res.ExecutedAt.Get().Format(time.RFC3339))
	}
	return diags
}

func resourceUserOffboardingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.LifecycleAPI.LifecycleUserOffboardingDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
