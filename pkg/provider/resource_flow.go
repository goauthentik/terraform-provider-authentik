package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceFlow() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
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
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedFlowDesignationEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedFlowDesignationEnumEnumValues),
			},
			"authentication": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHENTICATIONENUM_NONE,
				Description:      helpers.EnumToDescription(api.AllowedAuthenticationEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedAuthenticationEnumEnumValues),
			},
			"layout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.FLOWLAYOUTENUM_STACKED,
				Description:      helpers.EnumToDescription(api.AllowedFlowLayoutEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedFlowLayoutEnumEnumValues),
			},
			"background": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional URL to an image which will be used as the background during the flow.",
				Default:     "/static/dist/assets/images/flow_background.jpg",
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"denied_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.DENIEDACTIONENUM_MESSAGE_CONTINUE,
			},
			"compatibility_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceFlowSchemaToModel(d *schema.ResourceData) *api.FlowRequest {
	m := api.FlowRequest{
		Name:              d.Get("name").(string),
		Slug:              d.Get("slug").(string),
		Title:             d.Get("title").(string),
		CompatibilityMode: new(d.Get("compatibility_mode").(bool)),
		Designation:       api.FlowDesignationEnum(d.Get("designation").(string)),
		Authentication:    api.AuthenticationEnum(d.Get("authentication").(string)).Ptr(),
		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		Layout:            api.FlowLayoutEnum(d.Get("layout").(string)).Ptr(),
		DeniedAction:      api.DeniedActionEnum(d.Get("denied_action").(string)).Ptr(),
		Background:        helpers.GetP[string](d, "background"),
	}
	return &m
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsInstancesCreate(ctx).FlowRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	return resourceFlowRead(ctx, d, m)
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.FlowsApi.FlowsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "title", res.Title)
	helpers.SetWrapper(d, "designation", res.Designation)
	helpers.SetWrapper(d, "authentication", res.Authentication)
	helpers.SetWrapper(d, "denied_action", res.DeniedAction)
	helpers.SetWrapper(d, "layout", res.Layout)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "compatibility_mode", res.CompatibilityMode)
	helpers.SetWrapper(d, "background", res.Background)
	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsInstancesUpdate(ctx, d.Id()).FlowRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.FlowsApi.FlowsInstancesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
