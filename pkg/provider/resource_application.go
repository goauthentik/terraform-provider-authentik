package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
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
				Description:      helpers.EnumToDescription(api.AllowedPolicyEngineModeEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedPolicyEngineModeEnumValues),
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
		Provider:         *api.NewNullableInt32(helpers.GetIntP(d, ("protocol_provider"))),
		OpenInNewTab:     api.PtrBool(d.Get("open_in_new_tab").(bool)),
		PolicyEngineMode: api.PolicyEngineMode(d.Get("policy_engine_mode").(string)).Ptr(),
		Group:            helpers.GetP[string](d, "group"),
		MetaLaunchUrl:    helpers.GetP[string](d, "meta_launch_url"),
		MetaDescription:  helpers.GetP[string](d, "meta_description"),
		MetaPublisher:    helpers.GetP[string](d, "meta_publisher"),
	}

	m.BackchannelProviders = []int32{}
	for _, bp := range d.Get("backchannel_providers").([]interface{}) {
		m.BackchannelProviders = append(m.BackchannelProviders, int32(bp.(int)))
	}
	return &m
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsCreate(ctx).ApplicationRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
	}
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreApplicationsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "group", res.Group)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "open_in_new_tab", res.OpenInNewTab)
	helpers.SetWrapper(d, "protocol_provider", 0)
	if prov := res.Provider.Get(); prov != nil {
		helpers.SetWrapper(d, "protocol_provider", int(*prov))
	}
	helpers.SetWrapper(d, "meta_launch_url", res.MetaLaunchUrl)
	if res.MetaIcon.IsSet() {
		helpers.SetWrapper(d, "meta_icon", res.MetaIcon.Get())
	}
	helpers.SetWrapper(d, "meta_description", res.MetaDescription)
	helpers.SetWrapper(d, "meta_publisher", res.MetaPublisher)
	helpers.SetWrapper(d, "policy_engine_mode", res.PolicyEngineMode)
	helpers.SetWrapper(d, "backchannel_providers", res.BackchannelProviders)
	return diags
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceApplicationSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreApplicationsUpdate(ctx, d.Id()).ApplicationRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if i, iok := d.GetOk("meta_icon"); iok {
		hr, err := c.client.CoreApi.CoreApplicationsSetIconUrlCreate(ctx, res.Slug).FilePathRequest(api.FilePathRequest{
			Url: i.(string),
		}).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
	}

	d.SetId(res.Slug)
	return resourceApplicationRead(ctx, d, m)
}

func resourceApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreApplicationsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
