package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceBrand() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrandRead,
		Description: "System --- Get brands by domain",
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
			"branding_default_flow_background": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"branding_custom_css": {
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
			"web_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_application": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceBrandRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.CoreApi.CoreBrandsList(ctx)
	if s, ok := d.GetOk("domain"); ok {
		req = req.Domain(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching brands found")
	}
	f := res.Results[0]
	d.SetId(f.BrandUuid)
	helpers.SetWrapper(d, "domain", f.Domain)
	helpers.SetWrapper(d, "default", f.Default)
	helpers.SetWrapper(d, "branding_title", f.BrandingTitle)
	helpers.SetWrapper(d, "branding_logo", f.BrandingLogo)
	helpers.SetWrapper(d, "branding_favicon", f.BrandingFavicon)
	helpers.SetWrapper(d, "branding_default_flow_background", f.BrandingDefaultFlowBackground)
	helpers.SetWrapper(d, "branding_custom_css", f.BrandingCustomCss)
	helpers.SetWrapper(d, "flow_authentication", f.FlowAuthentication.Get())
	helpers.SetWrapper(d, "flow_invalidation", f.FlowInvalidation.Get())
	helpers.SetWrapper(d, "flow_recovery", f.FlowRecovery.Get())
	helpers.SetWrapper(d, "flow_unenrollment", f.FlowUnenrollment.Get())
	helpers.SetWrapper(d, "flow_user_settings", f.FlowUserSettings.Get())
	helpers.SetWrapper(d, "flow_device_code", f.FlowDeviceCode.Get())
	helpers.SetWrapper(d, "web_certificate", f.WebCertificate.Get())
	helpers.SetWrapper(d, "client_certificates", f.ClientCertificates)
	helpers.SetWrapper(d, "default_application", f.DefaultApplication.Get())
	return diags
}
