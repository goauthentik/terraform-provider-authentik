package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceSourceTelegram() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceTelegramCreate,
		ReadContext:   resourceSourceTelegramRead,
		UpdateContext: resourceSourceTelegramUpdate,
		DeleteContext: resourceSourceTelegramDelete,
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
			"pre_authentication_flow": {
				Type:     schema.TypeString,
				Required: true,
			},

			"bot_username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bot_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_message_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"property_mappings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"property_mappings_group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSourceTelegramSchemaToSource(d *schema.ResourceData) *api.TelegramSourceRequest {
	r := api.TelegramSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          new(d.Get("enabled").(bool)),
		UserPathTemplate: new(d.Get("user_path_template").(string)),
		PolicyEngineMode: api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		UserMatchingMode: api.UserMatchingModeEnum(d.Get("user_matching_mode").(string)).Ptr(),

		AuthenticationFlow:    *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
		EnrollmentFlow:        *api.NewNullableString(helpers.GetP[string](d, "enrollment_flow")),
		PreAuthenticationFlow: d.Get("pre_authentication_flow").(string),
		UserPropertyMappings:  helpers.CastSlice[string](d, "property_mappings"),
		GroupPropertyMappings: helpers.CastSlice[string](d, "property_mappings_group"),

		BotUsername:          d.Get("bot_username").(string),
		BotToken:             d.Get("bot_token").(string),
		RequestMessageAccess: helpers.GetP[bool](d, "request_message_access"),
	}
	return &r
}

func resourceSourceTelegramCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceTelegramSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesTelegramCreate(ctx).TelegramSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceTelegramRead(ctx, d, m)
}

func resourceSourceTelegramRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesTelegramRetrieve(ctx, d.Id()).Execute()
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
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "user_matching_mode", res.UserMatchingMode)
	helpers.SetWrapper(d, "pre_authentication_flow", res.PreAuthenticationFlow)

	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.UserPropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.GroupPropertyMappings,
	))
	helpers.SetWrapper(d, "bot_username", res.BotUsername)
	helpers.SetWrapper(d, "request_message_access", res.RequestMessageAccess)
	return diags
}

func resourceSourceTelegramUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceTelegramSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesTelegramUpdate(ctx, d.Id()).TelegramSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceTelegramRead(ctx, d, m)
}

func resourceSourceTelegramDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesTelegramDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
