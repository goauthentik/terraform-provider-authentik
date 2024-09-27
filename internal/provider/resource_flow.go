package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      EnumToDescription(api.AllowedFlowDesignationEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedFlowDesignationEnumEnumValues),
			},
			"authentication": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.AUTHENTICATIONENUM_NONE,
				Description:      EnumToDescription(api.AllowedAuthenticationEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedAuthenticationEnumEnumValues),
			},
			"layout": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.FLOWLAYOUTENUM_STACKED,
				Description:      EnumToDescription(api.AllowedFlowLayoutEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedFlowLayoutEnumEnumValues),
			},
			"background": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional URL to an image which will be used as the background during the flow.",
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPolicyEngineModeEnumValues),
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
		CompatibilityMode: api.PtrBool(d.Get("compatibility_mode").(bool)),
		Designation:       api.FlowDesignationEnum(d.Get("designation").(string)),
		Authentication:    api.AuthenticationEnum(d.Get("authentication").(string)).Ptr(),
		PolicyEngineMode:  api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		Layout:            api.FlowLayoutEnum(d.Get("layout").(string)).Ptr(),
		DeniedAction:      api.DeniedActionEnum(d.Get("denied_action").(string)).Ptr(),
	}
	return &m
}

func resourceFlowCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsInstancesCreate(ctx).FlowRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	if bg, ok := d.GetOk("background"); ok {
		hr, err := c.client.FlowsApi.FlowsInstancesSetBackgroundUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: bg.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.FlowsApi.FlowsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "title", res.Title)
	setWrapper(d, "designation", res.Designation)
	setWrapper(d, "authentication", res.Authentication)
	setWrapper(d, "denied_action", res.DeniedAction)
	setWrapper(d, "layout", res.Layout)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "compatibility_mode", res.CompatibilityMode)
	if _, bg := d.GetOk("background"); bg {
		setWrapper(d, "background", res.Background)
	}
	return diags
}

func resourceFlowUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceFlowSchemaToModel(d)

	res, hr, err := c.client.FlowsApi.FlowsInstancesUpdate(ctx, d.Id()).FlowRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	if bg, ok := d.GetOk("background"); ok {
		hr, err := c.client.FlowsApi.FlowsInstancesSetBackgroundUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: bg.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}
	return resourceFlowRead(ctx, d, m)
}

func resourceFlowDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.FlowsApi.FlowsInstancesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
