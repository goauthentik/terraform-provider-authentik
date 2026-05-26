package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

var (
	resourceOutpostDeleteRetryTimeout  = 10 * time.Second
	resourceOutpostDeleteRetryInterval = 250 * time.Millisecond
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
		Providers:         helpers.CastSliceInt32(d.Get("protocol_providers").([]any)),
	}

	if l, ok := d.Get("config").(string); ok && l != "" {
		attr, err := helpers.GetJSON[map[string]any](d, ("config"))
		m.Config = attr
		return &m, err
	} else {
		defaultConfig, hr, err := c.client.OutpostsAPI.OutpostsInstancesDefaultSettingsRetrieve(context.Background()).Execute()
		if err != nil {
			return nil, helpers.HTTPToDiag(d, hr, err)
		}
		m.Config = defaultConfig.Config
	}
	return &m, nil
}

func resourceOutpostCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceOutpostSchemaToModel(d, c)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsAPI.OutpostsInstancesCreate(ctx).OutpostRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsAPI.OutpostsInstancesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "type", res.Type)
	helpers.SetWrapper(d, "protocol_providers", helpers.ListConsistentMerge(
		helpers.CastSlice[int](d, "protocol_providers"),
		helpers.Slice32ToInt(res.Providers),
	))
	helpers.SetWrapper(d, "service_connection", res.ServiceConnection.Get())
	return helpers.SetJSON(d, "config", res.Config)
}

func resourceOutpostUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceOutpostSchemaToModel(d, c)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsAPI.OutpostsInstancesUpdate(ctx, d.Id()).OutpostRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceOutpostRead(ctx, d, m)
}

func resourceOutpostDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id := d.Id()
	hr, err := c.client.OutpostsAPI.OutpostsInstancesDestroy(ctx, id).Execute()
	if err != nil {
		if hr != nil && hr.StatusCode == http.StatusMethodNotAllowed {
			if resourceOutpostDeletedAfterMethodNotAllowed(ctx, c, id) {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return helpers.HTTPToDiag(d, hr, err)
	}
	d.SetId("")
	return diag.Diagnostics{}
}

func resourceOutpostDeletedAfterMethodNotAllowed(ctx context.Context, c *APIClient, id string) bool {
	if resourceOutpostRetrieveNotFound(ctx, c, id) {
		return true
	}

	timeout := time.NewTimer(resourceOutpostDeleteRetryTimeout)
	defer timeout.Stop()

	ticker := time.NewTicker(resourceOutpostDeleteRetryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-timeout.C:
			return false
		case <-ticker.C:
			if resourceOutpostRetrieveNotFound(ctx, c, id) {
				return true
			}
		}
	}
}

func resourceOutpostRetrieveNotFound(ctx context.Context, c *APIClient, id string) bool {
	_, hr, err := c.client.OutpostsAPI.OutpostsInstancesRetrieve(ctx, id).Execute()
	return err != nil && hr != nil && hr.StatusCode == http.StatusNotFound
}
