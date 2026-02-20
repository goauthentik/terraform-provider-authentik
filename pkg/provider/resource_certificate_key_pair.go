package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceCertificateKeyPair() *schema.Resource {
	return &schema.Resource{
		Description:   "System --- ",
		CreateContext: resourceCertificateKeyPairCreate,
		ReadContext:   resourceCertificateKeyPairRead,
		UpdateContext: resourceCertificateKeyPairUpdate,
		DeleteContext: resourceCertificateKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_data": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCertificateKeyPairSchemaToModel(d *schema.ResourceData) *api.CertificateKeyPairRequest {
	app := api.CertificateKeyPairRequest{
		Name:            d.Get("name").(string),
		CertificateData: d.Get("certificate_data").(string),
		KeyData:         helpers.GetP[string](d, "key_data"),
	}
	return &app
}

func resourceCertificateKeyPairCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceCertificateKeyPairSchemaToModel(d)

	res, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsCreate(ctx).CertificateKeyPairRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceCertificateKeyPairRead(ctx, d, m)
}

func resourceCertificateKeyPairRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)

	rc, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsViewCertificateRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "certificate_data", rc.Data+"\n")

	rk, _, err := c.client.CryptoApi.CryptoCertificatekeypairsViewPrivateKeyRetrieve(ctx, d.Id()).Execute()
	if err == nil {
		helpers.SetWrapper(d, "key_data", rk.Data+"\n")
	}

	return diags
}

func resourceCertificateKeyPairUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceCertificateKeyPairSchemaToModel(d)

	res, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsUpdate(ctx, d.Id()).CertificateKeyPairRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceCertificateKeyPairRead(ctx, d, m)
}

func resourceCertificateKeyPairDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CryptoApi.CryptoCertificatekeypairsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
