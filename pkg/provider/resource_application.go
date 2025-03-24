package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceApplication() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
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
			"protocol_provider": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backchannel_providers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"meta_launch_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_icon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"meta_publisher": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_engine_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.POLICYENGINEMODE_ANY,
				Description:      EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedPolicyEngineModeEnumValues),
			},
			"open_in_new_tab": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceApplicationSchemaToModel(d *schema.ResourceData) *api.ApplicationRequest {
	m := api.ApplicationRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Provider:         api.NullableInt32{},
		OpenInNewTab:     api.PtrBool(d.Get("open_in_new_tab").(bool)),
		PolicyEngineMode: api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
	}

	if p, pSet := d.GetOk("protocol_provider"); pSet {
		m.Provider.Set(api.PtrInt32(int32(p.(int))))
	} else {
		m.Provider.Set(nil)
	}
	m.BackchannelProviders = []int32{}
	for _, bp := range d.Get("backchannel_providers").([]interface{}) {
		m.BackchannelProviders = append(m.BackchannelProviders, int32(bp.(int)))
	}

	if l, ok := d.Get("group").(string); ok {
		m.Group = &l
	}
	if l, ok := d.Get("meta_launch_url").(string); ok {
		m.MetaLaunchUrl = &l
	}
	if l, ok := d.Get("meta_description").(string); ok {
		m.MetaDescription = &l
	}
	if l, ok := d.Get("meta_publisher").(string); ok {
		m.MetaPublisher = &l
	}
	return &m
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsCreate(ctx).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreApplicationsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "name", res.Name)
	setWrapper(d, "group", res.Group)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "open_in_new_tab", res.OpenInNewTab)
	setWrapper(d, "protocol_provider", 0)
	if prov := res.Provider.Get(); prov != nil {
		setWrapper(d, "protocol_provider", int(*prov))
	}
	setWrapper(d, "meta_launch_url", res.MetaLaunchUrl)
	if res.MetaIcon.IsSet() {
		setWrapper(d, "meta_icon", res.MetaIcon.Get())
	}
	setWrapper(d, "meta_description", res.MetaDescription)
	setWrapper(d, "meta_publisher", res.MetaPublisher)
	setWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	setWrapper(d, "backchannel_providers", res.BackchannelProviders)
	return diags
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsUpdate(ctx, d.Id()).ApplicationRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return httpToDiag(d, hr, err)
		}
	}

	d.SetId(res.Slug)
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreApplicationsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
