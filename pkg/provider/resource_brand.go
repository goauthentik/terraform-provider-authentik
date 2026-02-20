package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
		},
	}
}

func resourceBrandSchemaToModel(d *schema.ResourceData) (*api.BrandRequest, diag.Diagnostics) {
	m := api.BrandRequest{
		Domain:                        d.Get("domain").(string),
		Default:                       new(d.Get("default").(bool)),
		BrandingTitle:                 helpers.GetP[string](d, "branding_title"),
		BrandingLogo:                  helpers.GetP[string](d, "branding_logo"),
		BrandingFavicon:               helpers.GetP[string](d, "branding_favicon"),
		BrandingDefaultFlowBackground: helpers.GetP[string](d, "branding_default_flow_background"),
		BrandingCustomCss:             helpers.GetP[string](d, "branding_custom_css"),
		FlowAuthentication:            *api.NewNullableString(helpers.GetP[string](d, "flow_authentication")),
		FlowInvalidation:              *api.NewNullableString(helpers.GetP[string](d, "flow_invalidation")),
		FlowRecovery:                  *api.NewNullableString(helpers.GetP[string](d, "flow_recovery")),
		FlowUnenrollment:              *api.NewNullableString(helpers.GetP[string](d, "flow_unenrollment")),
		FlowUserSettings:              *api.NewNullableString(helpers.GetP[string](d, "flow_user_settings")),
		FlowDeviceCode:                *api.NewNullableString(helpers.GetP[string](d, "flow_device_code")),
		WebCertificate:                *api.NewNullableString(helpers.GetP[string](d, "web_certificate")),
		ClientCertificates:            helpers.CastSlice[string](d, "client_certificates"),
		DefaultApplication:            *api.NewNullableString(helpers.GetP[string](d, "default_application")),
	}

	attr, err := helpers.GetJSON[map[string]any](d, ("attributes"))
	m.Attributes = attr
	return &m, err
}

func resourceBrandCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	mo, diags := resourceBrandSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreBrandsCreate(ctx).BrandRequest(*mo).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.BrandUuid)
	return resourceBrandRead(ctx, d, m)
}

func resourceBrandRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreBrandsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "domain", res.Domain)
	helpers.SetWrapper(d, "branding_title", res.BrandingTitle)
	helpers.SetWrapper(d, "branding_logo", res.BrandingLogo)
	helpers.SetWrapper(d, "branding_favicon", res.BrandingFavicon)
	helpers.SetWrapper(d, "branding_default_flow_background", res.BrandingDefaultFlowBackground)
	helpers.SetWrapper(d, "branding_custom_css", res.BrandingCustomCss)
	helpers.SetWrapper(d, "flow_authentication", res.FlowAuthentication.Get())
	helpers.SetWrapper(d, "flow_invalidation", res.FlowInvalidation.Get())
	helpers.SetWrapper(d, "flow_recovery", res.FlowRecovery.Get())
	helpers.SetWrapper(d, "flow_unenrollment", res.FlowUnenrollment.Get())
	helpers.SetWrapper(d, "flow_user_settings", res.FlowUserSettings.Get())
	helpers.SetWrapper(d, "flow_device_code", res.FlowDeviceCode.Get())
	helpers.SetWrapper(d, "web_certificate", res.WebCertificate.Get())
	helpers.SetWrapper(d, "client_certificates", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "client_certificates"),
		res.ClientCertificates,
	))
	helpers.SetWrapper(d, "default_application", res.DefaultApplication.Get())
	return helpers.SetJSON(d, "attributes", res.Attributes)
}

func resourceBrandUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	obj, diags := resourceBrandSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreBrandsUpdate(ctx, d.Id()).BrandRequest(*obj).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.BrandUuid)
	return resourceBrandRead(ctx, d, m)
}

func resourceBrandDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreBrandsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
