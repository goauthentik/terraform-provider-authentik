package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Default:          api.MUTUALTLSSTAGEMODEENUM_OPTIONAL,
				Description:      EnumToDescription(api.AllowedMutualTLSStageModeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedMutualTLSStageModeEnumEnumValues),
			},
			"cert_attribute": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.CERTATTRIBUTEENUM_EMAIL,
				Description:      EnumToDescription(api.AllowedCertAttributeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedCertAttributeEnumEnumValues),
			},
			"user_attribute": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.USERATTRIBUTEENUM_EMAIL,
				Description:      EnumToDescription(api.AllowedUserAttributeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserAttributeEnumEnumValues),
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
		Mode:                   api.MutualTLSStageModeEnum(d.Get("mode").(string)),
		CertAttribute:          api.CertAttributeEnum(d.Get("cert_attribute").(string)),
		UserAttribute:          api.UserAttributeEnum(d.Get("user_attribute").(string)),
		CertificateAuthorities: castSlice[string](d.Get("certificate_authorities").([]interface{})),
	}
	return &r
}

func resourceStageMutualTLSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStageMutualTLSSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesMtlsCreate(ctx).MutualTLSStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageMutualTLSRead(ctx, d, m)
}

func resourceStageMutualTLSRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesMtlsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "mode", res.Mode)
	setWrapper(d, "cert_attribute", res.CertAttribute)
	setWrapper(d, "user_attribute", res.UserAttribute)
	localCertificateAuthorities := castSlice[string](d.Get("certificate_authorities").([]interface{}))
	setWrapper(d, "certificate_authorities", listConsistentMerge(localCertificateAuthorities, res.CertificateAuthorities))
	return diags
}

func resourceStageMutualTLSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStageMutualTLSSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesMtlsUpdate(ctx, d.Id()).MutualTLSStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStageMutualTLSRead(ctx, d, m)
}

func resourceStageMutualTLSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesMtlsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
