package provider

import (
	"context"
	"encoding/json"

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
			"branding_default_flow_background": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/static/dist/assets/images/flow_background.jpg",
			},
			"branding_custom_css": {
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
			},
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      JSONDescription,
				DiffSuppressFunc: diffSuppressJSON,
				ValidateDiagFunc: ValidateJSON,
			},
		},
	}
}

func resourceBrandSchemaToModel(d *schema.ResourceData) (*api.BrandRequest, diag.Diagnostics) {
	m := api.BrandRequest{
		ClientCertificates:            castSlice[string](d.Get("client_certificates").([]interface{})),
		Domain:                        d.Get("domain").(string),
		Default:                       api.PtrBool(d.Get("default").(bool)),
		BrandingTitle:                 getP[string](d, "branding_title"),
		BrandingLogo:                  getP[string](d, "branding_logo"),
		BrandingFavicon:               getP[string](d, "branding_favicon"),
		BrandingDefaultFlowBackground: getP[string](d, "branding_default_flow_background"),
		BrandingCustomCss:             getP[string](d, "branding_custom_css"),
		FlowAuthentication:            *api.NewNullableString(getP[string](d, "flow_authentication")),
		FlowInvalidation:              *api.NewNullableString(getP[string](d, "flow_invalidation")),
		FlowRecovery:                  *api.NewNullableString(getP[string](d, "flow_recovery")),
		FlowUnenrollment:              *api.NewNullableString(getP[string](d, "flow_unenrollment")),
		FlowUserSettings:              *api.NewNullableString(getP[string](d, "flow_user_settings")),
		FlowDeviceCode:                *api.NewNullableString(getP[string](d, "flow_device_code")),
		WebCertificate:                *api.NewNullableString(getP[string](d, "web_certificate")),
		DefaultApplication:            *api.NewNullableString(getP[string](d, "default_application")),
	}

	attr, err := getJSON[map[string]interface{}](d, ("attributes"))
	m.Attributes = attr
	return &m, err
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
	setWrapper(d, "branding_default_flow_background", res.BrandingDefaultFlowBackground)
	setWrapper(d, "branding_custom_css", res.BrandingCustomCss)
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
	if res.HasClientCertificates() {
		localClientCertificates := castSlice[string](d.Get("client_certificates").([]interface{}))
		setWrapper(d, "client_certificates", listConsistentMerge(localClientCertificates, res.ClientCertificates))
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
