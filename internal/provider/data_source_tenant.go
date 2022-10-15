package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTenant() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTenantRead,
		Description: "Get tenants by domain",
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"branding_title": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"branding_logo": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"branding_favicon": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_invalidation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_recovery": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_unenrollment": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_user_settings": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flow_device_code": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"event_retention": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"web_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.CoreApi.CoreTenantsList(ctx)
	if s, ok := d.GetOk("domain"); ok {
		req = req.Domain(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching tenants found")
	}
	f := res.Results[0]
	d.SetId(f.TenantUuid)
	setWrapper(d, "domain", f.Domain)
	setWrapper(d, "default", f.Default)
	setWrapper(d, "branding_title", f.BrandingTitle)
	setWrapper(d, "branding_logo", f.BrandingLogo)
	setWrapper(d, "branding_favicon", f.BrandingFavicon)
	setWrapper(d, "flow_authentication", f.FlowAuthentication.Get())
	setWrapper(d, "flow_invalidation", f.FlowInvalidation.Get())
	setWrapper(d, "flow_recovery", f.FlowRecovery.Get())
	setWrapper(d, "flow_unenrollment", f.FlowUnenrollment.Get())
	setWrapper(d, "flow_user_settings", f.FlowUserSettings.Get())
	setWrapper(d, "flow_device_code", f.FlowDeviceCode.Get())
	setWrapper(d, "event_retention", f.EventRetention)
	setWrapper(d, "web_certificate", f.WebCertificate.Get())
	return diags
}
