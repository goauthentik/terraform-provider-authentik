package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceStageMutualTLS() *schema.Resource {
	return &schema.Resource{
		Description:   "Flows & Stages --- ",
		CreateContext: resourceStageMutualTLSCreate,
		ReadContext:   resourceStageMutualTLSRead,
		UpdateContext: resourceStageMutualTLSUpdate,
		DeleteContext: resourceStageMutualTLSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.STAGEMODEENUM_OPTIONAL,
				Description:      helpers.EnumToDescription(api.AllowedStageModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedStageModeEnumEnumValues),
			},
			"cert_attribute": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CERTATTRIBUTEENUM_EMAIL,
				Description:      helpers.EnumToDescription(api.AllowedCertAttributeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedCertAttributeEnumEnumValues),
			},
			"user_attribute": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERATTRIBUTEENUM_EMAIL,
				Description:      helpers.EnumToDescription(api.AllowedUserAttributeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserAttributeEnumEnumValues),
			},
			"certificate_authorities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceStageMutualTLSSchemaToProvider(d *schema.ResourceData) *api.MutualTLSStageRequest {
	r := api.MutualTLSStageRequest{
		Name:                   d.Get("name").(string),
		Mode:                   api.StageModeEnum(d.Get("mode").(string)),
		CertAttribute:          api.CertAttributeEnum(d.Get("cert_attribute").(string)),
		UserAttribute:          api.UserAttributeEnum(d.Get("user_attribute").(string)),
		CertificateAuthorities: helpers.CastSlice[string](d, "certificate_authorities"),
	}
	return &r
}

func resourceStageMutualTLSCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageMutualTLSSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesMtlsCreate(ctx).MutualTLSStageRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageMutualTLSRead(ctx, d, m)
}

func resourceStageMutualTLSRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesMtlsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "mode", res.Mode)
	helpers.SetWrapper(d, "cert_attribute", res.CertAttribute)
	helpers.SetWrapper(d, "user_attribute", res.UserAttribute)
	helpers.SetWrapper(d, "certificate_authorities", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "certificate_authorities"),
		res.CertificateAuthorities,
	))
	return diags
}

func resourceStageMutualTLSUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageMutualTLSSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesMtlsUpdate(ctx, d.Id()).MutualTLSStageRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageMutualTLSRead(ctx, d, m)
}

func resourceStageMutualTLSDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesMtlsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
