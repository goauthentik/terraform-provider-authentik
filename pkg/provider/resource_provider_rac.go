package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
)

func resourceProviderRAC() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderRACCreate,
		ReadContext:   resourceProviderRACRead,
		UpdateContext: resourceProviderRACUpdate,
		DeleteContext: resourceProviderRACDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"authentication_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"authorization_flow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"connection_expiry": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "seconds=0",
				Description:      helpers.RelativeDurationDescription,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
			},
		},
	}
}

func resourceProviderRACSchemaToProvider(d *schema.ResourceData) (*api.RACProviderRequest, diag.Diagnostics) {
	r := api.RACProviderRequest{
		Name:               d.Get("name").(string),
		AuthorizationFlow:  d.Get("authorization_flow").(string),
		PropertyMappings:   helpers.CastSlice[string](d.Get("property_mappings").([]interface{})),
		ConnectionExpiry:   api.PtrString(d.Get("connection_expiry").(string)),
		AuthenticationFlow: *api.NewNullableString(helpers.GetP[string](d, "authentication_flow")),
	}
	attr, err := helpers.GetJSON[map[string]interface{}](d, ("settings"))
	r.Settings = attr
	return &r, err
}

func resourceProviderRACCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersRacCreate(ctx).RACProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRACRead(ctx, d, m)
}

func resourceProviderRACRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersRacRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "authentication_flow", res.AuthenticationFlow.Get())
	helpers.SetWrapper(d, "authorization_flow", res.AuthorizationFlow)
	helpers.SetWrapper(d, "connection_expiry", res.ConnectionExpiry)
	localMappings := helpers.CastSlice[string](d.Get("property_mappings").([]interface{}))
	if len(localMappings) > 0 {
		helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(localMappings, res.PropertyMappings))
	}
	b, err := json.Marshal(res.Settings)
	if err != nil {
		return diag.FromErr(err)
	}
	helpers.SetWrapper(d, "settings", string(b))
	return diags
}

func resourceProviderRACUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app, diags := resourceProviderRACSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersRacUpdate(ctx, int32(id)).RACProviderRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderRACRead(ctx, d, m)
}

func resourceProviderRACDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersRacDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
