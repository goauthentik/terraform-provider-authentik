package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceSourcePlex() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourcePlexCreate,
		ReadContext:   resourceSourcePlexRead,
		UpdateContext: resourceSourcePlexUpdate,
		DeleteContext: resourceSourcePlexDelete,
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
				Optional: true,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_path_template": {
				Type:     schema.TypeString,
				Default:  "goauthentik.io/sources/%(slug)s",
				Optional: true,
			},
			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enrollment_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"promoted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"user_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedUserMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserMatchingModeEnumEnumValues),
			},
			"group_matching_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.GROUPMATCHINGMODEENUM_IDENTIFIER,
				Description:      helpers.EnumToDescription(api.AllowedGroupMatchingModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedGroupMatchingModeEnumEnumValues),
			},

			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allow_friends": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"plex_token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSourcePlexSchemaToSource(d *schema.ResourceData) *api.PlexSourceRequest {
	r := api.PlexSourceRequest{
		Name:               d.Get("name").(string),
		Slug:               d.Get("slug").(string),
		Enabled:            new(d.Get("enabled").(bool)),
		Promoted:           new(d.Get("promoted").(bool)),
		UserPathTemplate:   new(d.Get("user_path_template").(string)),
		PolicyEngineMode:   api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode:   api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),
		GroupMatchingMode:  api.GroupMatchingModeEnum(d.Get("group_matching_mode").(string)).Ptr(),
		AuthenticationFlow: *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		EnrollmentFlow:     *api.NewNullableString(helpers.GetP[string](d, "enrollment_flow")),

		ClientId:       new(d.Get("client_id").(string)),
		AllowFriends:   new(d.Get("allow_friends").(bool)),
		PlexToken:      d.Get("plex_token").(string),
		AllowedServers: helpers.CastSlice[string](d, "allowed_servers"),
	}
	return &r
}

func resourceSourcePlexCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourcePlexSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesPlexCreate(ctx).PlexSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesPlexRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "user_path_template", res.UserPathTemplate)

	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "enrollment_flow", res.EnrollmentFlow.Get())
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "promoted", res.Promoted)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "user_matching_mode", res.UserMatchingMode)
	helpers.SetWrapper(d, "group_matching_mode", res.GroupMatchingMode)

	helpers.SetWrapper(d, "client_id", res.ClientId)
	helpers.SetWrapper(d, "allowed_servers", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "allowed_servers"),
		res.AllowedServers,
	))
	helpers.SetWrapper(d, "allow_friends", res.AllowFriends)
	helpers.SetWrapper(d, "plex_token", res.PlexToken)
	return diags
}

func resourceSourcePlexUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourcePlexSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesPlexUpdate(ctx, d.Id()).PlexSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourcePlexRead(ctx, d, m)
}

func resourceSourcePlexDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesPlexDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
