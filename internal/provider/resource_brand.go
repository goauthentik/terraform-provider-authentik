package provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceBrand() *schema.Resource {
	return &schema.Resource{
		Description:   "System --- ",
		CreateContext: resourceBrandCreate,
		ReadContext:   resourceBrandRead,
		UpdateContext: resourceBrandUpdate,
		DeleteContext: resourceBrandDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"branding_title": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "authentik",
			},
			"branding_logo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"branding_favicon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_authentication": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_invalidation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_recovery": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_unenrollment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_user_settings": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flow_device_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_application": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      "JSON format expected. Use jsonencode() to pass objects.",
				DiffSuppressFunc: diffSuppressJSON,
			},
		},
	}
}

func resourceBrandSchemaToModel(d *schema.ResourceData) (*api.BrandRequest, diag.Diagnostics) {
	m := api.BrandRequest{
		Domain:  d.Get("domain").(string),
		Default: api.PtrBool(d.Get("default").(bool)),
	}

	if l, ok := d.Get("branding_title").(string); ok {
		m.BrandingTitle = &l
	}
	if l, ok := d.Get("branding_logo").(string); ok {
		m.BrandingLogo = &l
	}
	if l, ok := d.Get("branding_favicon").(string); ok {
		m.BrandingFavicon = &l
	}

	if l, ok := d.Get("flow_authentication").(string); ok {
		m.FlowAuthentication.Set(&l)
	} else {
		m.FlowAuthentication.Set(nil)
	}

	if l, ok := d.Get("flow_invalidation").(string); ok {
		m.FlowInvalidation.Set(&l)
	} else {
		m.FlowInvalidation.Set(nil)
	}

	if l, ok := d.Get("flow_recovery").(string); ok {
		m.FlowRecovery.Set(&l)
	} else {
		m.FlowRecovery.Set(nil)
	}

	if l, ok := d.Get("flow_unenrollment").(string); ok {
		m.FlowUnenrollment.Set(&l)
	} else {
		m.FlowUnenrollment.Set(nil)
	}

	if l, ok := d.Get("flow_user_settings").(string); ok {
		m.FlowUserSettings.Set(&l)
	} else {
		m.FlowUserSettings.Set(nil)
	}

	if l, ok := d.Get("flow_device_code").(string); ok {
		m.FlowDeviceCode.Set(&l)
	} else {
		m.FlowDeviceCode.Set(nil)
	}

	if l, ok := d.Get("web_certificate").(string); ok {
		m.WebCertificate.Set(&l)
	} else {
		m.WebCertificate.Set(nil)
	}

	if l, ok := d.Get("default_application").(string); ok {
		m.DefaultApplication.Set(&l)
	} else {
		m.DefaultApplication.Set(nil)
	}

	attr := make(map[string]interface{})
	if l, ok := d.Get("attributes").(string); ok && l != "" {
		err := json.NewDecoder(strings.NewReader(l)).Decode(&attr)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}
	m.Attributes = attr
	return &m, nil
}

func resourceBrandCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mo, diags := resourceBrandSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreBrandsCreate(ctx).BrandRequest(*mo).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.BrandUuid)
	return resourceBrandRead(ctx, d, m)
}

func resourceBrandRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreBrandsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "domain", res.Domain)
	setWrapper(d, "branding_title", res.BrandingTitle)
	setWrapper(d, "branding_logo", res.BrandingLogo)
	setWrapper(d, "branding_favicon", res.BrandingFavicon)
	if res.FlowAuthentication.IsSet() {
		setWrapper(d, "flow_authentication", res.FlowAuthentication.Get())
	}
	if res.FlowInvalidation.IsSet() {
		setWrapper(d, "flow_invalidation", res.FlowInvalidation.Get())
	}
	if res.FlowRecovery.IsSet() {
		setWrapper(d, "flow_recovery", res.FlowRecovery.Get())
	}
	if res.FlowUnenrollment.IsSet() {
		setWrapper(d, "flow_unenrollment", res.FlowUnenrollment.Get())
	}
	if res.FlowUserSettings.IsSet() {
		setWrapper(d, "flow_user_settings", res.FlowUserSettings.Get())
	}
	if res.FlowDeviceCode.IsSet() {
		setWrapper(d, "flow_device_code", res.FlowDeviceCode.Get())
	}
	if res.WebCertificate.IsSet() {
		setWrapper(d, "web_certificate", res.WebCertificate.Get())
	}
	if res.DefaultApplication.IsSet() {
		setWrapper(d, "default_application", res.DefaultApplication.Get())
	}
	b, err := json.Marshal(res.Attributes)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "attributes", string(b))
	return diags
}

func resourceBrandUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	obj, diags := resourceBrandSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreBrandsUpdate(ctx, d.Id()).BrandRequest(*obj).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.BrandUuid)
	return resourceBrandRead(ctx, d, m)
}

func resourceBrandDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreBrandsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
