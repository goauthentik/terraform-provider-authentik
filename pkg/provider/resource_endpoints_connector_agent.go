package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceEndpointsConnectorAgent() *schema.Resource {
	return &schema.Resource{
		Description:   "Endpoint Devices --- ",
		CreateContext: resourceEndpointsConnectorAgentCreate,
		ReadContext:   resourceEndpointsConnectorAgentRead,
		UpdateContext: resourceEndpointsConnectorAgentUpdate,
		DeleteContext: resourceEndpointsConnectorAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"snapshot_expiry": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "hours=24",
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"auth_session_duration": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "hours=8",
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"auth_terminate_session_on_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"refresh_interval": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "minutes=30",
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nss_uid_offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2000,
			},
			"nss_gid_offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4000,
			},
			"challenge_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"challenge_idle_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "seconds=3",
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"challenge_trigger_check_in": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"jwt_federation_providers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceEndpointsConnectorAgentSchemaToProvider(d *schema.ResourceData) (*api.AgentConnectorRequest, diag.Diagnostics) {
	r := api.AgentConnectorRequest{
		Name:                         d.Get("name").(string),
		Enabled:                      helpers.GetP[bool](d, "enabled"),
		SnapshotExpiry:               helpers.GetP[string](d, "snapshot_expiry"),
		AuthSessionDuration:          helpers.GetP[string](d, "auth_session_duration"),
		AuthTerminateSessionOnExpiry: helpers.GetP[bool](d, "auth_terminate_session_on_expiry"),
		RefreshInterval:              helpers.GetP[string](d, "refresh_interval"),
		AuthorizationFlow:            *api.NewNullableString(helpers.GetP[string](d, "authorization_flow")),
		NssUidOffset:                 helpers.GetIntP(d, "nss_uid_offset"),
		NssGidOffset:                 helpers.GetIntP(d, "nss_gid_offset"),
		ChallengeKey:                 *api.NewNullableString(helpers.GetP[string](d, "challenge_key")),
		ChallengeIdleTimeout:         helpers.GetP[string](d, "challenge_idle_timeout"),
		ChallengeTriggerCheckIn:      helpers.GetP[bool](d, "challenge_trigger_check_in"),
		JwtFederationProviders:       helpers.CastSliceInt32(d.Get("jwt_federation_providers").([]any)),
	}

	return &r, nil
}

func resourceEndpointsConnectorAgentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceEndpointsConnectorAgentSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsApi.EndpointsAgentsConnectorsCreate(ctx).AgentConnectorRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(*res.ConnectorUuid)
	return resourceEndpointsConnectorAgentRead(ctx, d, m)
}

func resourceEndpointsConnectorAgentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.EndpointsApi.EndpointsAgentsConnectorsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "snapshot_expiry", res.SnapshotExpiry)
	helpers.SetWrapper(d, "auth_session_duration", res.AuthSessionDuration)
	helpers.SetWrapper(d, "auth_terminate_session_on_expiry", res.AuthTerminateSessionOnExpiry)
	helpers.SetWrapper(d, "refresh_interval", res.RefreshInterval)
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow.Get())
	helpers.SetWrapper(d, "nss_uid_offset", res.NssUidOffset)
	helpers.SetWrapper(d, "nss_gid_offset", res.NssGidOffset)
	helpers.SetWrapper(d, "challenge_key", res.ChallengeKey.Get())
	helpers.SetWrapper(d, "challenge_idle_timeout", res.ChallengeIdleTimeout)
	helpers.SetWrapper(d, "challenge_trigger_check_in", res.ChallengeTriggerCheckIn)
	helpers.SetWrapper(d, "jwt_federation_providers", helpers.ListConsistentMerge(
		helpers.CastSlice[int](d, "jwt_federation_providers"),
		helpers.Slice32ToInt(res.JwtFederationProviders),
	))
	return diags
}

func resourceEndpointsConnectorAgentUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceEndpointsConnectorAgentSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.EndpointsApi.EndpointsAgentsConnectorsUpdate(ctx, d.Id()).AgentConnectorRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(*res.ConnectorUuid)
	return resourceEndpointsConnectorAgentRead(ctx, d, m)
}

func resourceEndpointsConnectorAgentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.EndpointsApi.EndpointsAgentsConnectorsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
