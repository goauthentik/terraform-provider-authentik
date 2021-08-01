package provider

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"branding_title": {
				Type:     schema.TypeBool,
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
		},
	}
}

func resourceTenantSchemaToModel(d *schema.ResourceData) (*api.TenantRequest, diag.Diagnostics) {
	m := api.TenantRequest{
		Domain: d.Get("domain").(string),
	}

	m.BrandingTitle = stringToPointer(d.Get("branding_title").(string))
	m.BrandingLogo = stringToPointer(d.Get("branding_logo").(string))
	m.BrandingFavicon = stringToPointer(d.Get("branding_favicon").(string))

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

	return &m, nil
}

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	mo, diags := resourceTenantSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreTenantsCreate(ctx).TenantRequest(*mo).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Domain)
	return resourceTenantRead(ctx, d, m)
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreTenantsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("domain", res.Domain)
	d.Set("branding_title", res.BrandingTitle)
	d.Set("branding_logo", res.BrandingLogo)
	d.Set("branding_favicon", res.BrandingFavicon)
	d.Set("flow_authentication", res.FlowAuthentication)
	d.Set("flow_invalidation", res.FlowInvalidation)
	d.Set("flow_recovery", res.FlowRecovery)
	d.Set("flow_unenrollment", res.FlowUnenrollment)
	return diags
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceTenantSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.CoreApi.CoreTenantsUpdate(ctx, d.Id()).TenantRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Domain)
	return resourceTenantRead(ctx, d, m)
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreTenantsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
