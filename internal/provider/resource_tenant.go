package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func resourceTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantCreate,
		ReadContext:   resourceTenantRead,
		UpdateContext: resourceTenantUpdate,
		DeleteContext: resourceTenantDelete,
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
			"event_retention": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "days=365",
			},
			"web_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTenantSchemaToModel(d *schema.ResourceData) *api.TenantRequest {
	m := api.TenantRequest{
		Domain:  d.Get("domain").(string),
		Default: boolToPointer(d.Get("default").(bool)),
	}

	m.EventRetention = stringToPointer(d.Get("event_retention").(string))

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

	if l, ok := d.Get("web_certificate").(string); ok {
		m.WebCertificate.Set(&l)
	} else {
		m.WebCertificate.Set(nil)
	}

	return &m
}

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mo := resourceTenantSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreTenantsCreate(ctx).TenantRequest(*mo).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.TenantUuid)
	return resourceTenantRead(ctx, d, m)
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreTenantsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("domain", res.Domain)
	d.Set("branding_title", res.BrandingTitle)
	d.Set("branding_logo", res.BrandingLogo)
	d.Set("branding_favicon", res.BrandingFavicon)
	if res.FlowAuthentication.IsSet() {
		d.Set("flow_authentication", res.FlowAuthentication.Get())
	}
	if res.FlowInvalidation.IsSet() {
		d.Set("flow_invalidation", res.FlowInvalidation.Get())
	}
	if res.FlowRecovery.IsSet() {
		d.Set("flow_recovery", res.FlowRecovery.Get())
	}
	if res.FlowUnenrollment.IsSet() {
		d.Set("flow_unenrollment", res.FlowUnenrollment.Get())
	}
	d.Set("event_retention", res.EventRetention)
	if res.WebCertificate.IsSet() {
		d.Set("web_certificate", res.WebCertificate.Get())
	}
	return diags
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceTenantSchemaToModel(d)

	res, hr, err := c.client.CoreApi.CoreTenantsUpdate(ctx, d.Id()).TenantRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.TenantUuid)
	return resourceTenantRead(ctx, d, m)
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreTenantsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
