package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
)

func resourceOutpost() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceOutpostCreate,
		ReadContext:   resourceOutpostRead,
		UpdateContext: resourceOutpostUpdate,
		DeleteContext: resourceOutpostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.OUTPOSTTYPEENUM_PROXY,
				Description:      helpers.EnumToDescription(api.AllowedOutpostTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedOutpostTypeEnumEnumValues),
			},
			"protocol_providers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"service_connection": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourceOutpostSchemaToModel(d *schema.ResourceData, c *APIClient) (*api.OutpostRequest, diag.Diagnostics) {
	m := api.OutpostRequest{
		Name:              d.Get("name").(string),
		Type:              api.OutpostTypeEnum(d.Get("type").(string)),
		ServiceConnection: *api.NewNullableString(helpers.GetP[string](d, "service_connection")),
		Providers:         helpers.CastSliceInt32(d.Get("protocol_providers").([]interface{})),
	}

	defaultConfig, hr, err := c.client.OutpostsApi.OutpostsInstancesDefaultSettingsRetrieve(context.Background()).Execute()
	if err != nil {
		return nil, helpers.HTTPToDiag(d, hr, err)
	}
	if l, ok := d.Get("config").(string); ok && l != "" {
		var c map[string]interface{}
		err := json.NewDecoder(strings.NewReader(l)).Decode(&c)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		m.Config = c
	} else {
		m.Config = defaultConfig.Config
	}
	return &m, nil
}

func resourceOutpostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceOutpostSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesCreate(ctx).OutpostRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "type", res.Type)
	localProviders := helpers.CastSlice[int](d.Get("protocol_providers").([]interface{}))
	helpers.SetWrapper(d, "protocol_providers", helpers.ListConsistentMerge(localProviders, helpers.Slice32ToInt(res.Providers)))
	if res.ServiceConnection.IsSet() {
		helpers.SetWrapper(d, "service_connection", res.ServiceConnection.Get())
	}
	b, err := json.Marshal(res.Config)
	if err != nil {
		return diag.FromErr(err)
	}
	helpers.SetWrapper(d, "config", string(b))
	return diags
}

func resourceOutpostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceOutpostSchemaToModel(d, c)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsApi.OutpostsInstancesUpdate(ctx, d.Id()).OutpostRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.OutpostsApi.OutpostsInstancesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
