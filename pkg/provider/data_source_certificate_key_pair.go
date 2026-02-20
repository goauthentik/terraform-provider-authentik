package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceCertificateKeyPair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateKeyPairRead,
		Description: "System --- Get certificate-key pairs by name",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fetch_certificate": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If set to true, certificate data will be fetched.",
			},
			"certificate_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA1-hashed certificate fingerprint",
			},
			"fingerprint256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA256-hashed certificate fingerprint",
			},
			"fetch_key": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If set to true, private key data will be fetched.",
			},
			"key_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceCertificateKeyPairRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.CryptoApi.CryptoCertificatekeypairsList(ctx)
	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	if len(res.Results) < 1 {
		return diag.Errorf("No matching groups found")
	}
	f := res.Results[0]

	d.SetId(f.Pk)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetWrapper(d, "expiry", f.CertExpiry.Get().String())
	helpers.SetWrapper(d, "subject", f.CertSubject.Get())
	helpers.SetWrapper(d, "fingerprint1", f.FingerprintSha1.Get())
	helpers.SetWrapper(d, "fingerprint256", f.FingerprintSha256.Get())

	if d.Get("fetch_certificate").(bool) {
		rc, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsViewCertificateRetrieve(ctx, d.Id()).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		helpers.SetWrapper(d, "certificate_data", rc.Data+"\n")
	}

	if d.Get("fetch_key").(bool) {
		rk, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsViewPrivateKeyRetrieve(ctx, d.Id()).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		helpers.SetWrapper(d, "key_data", rk.Data+"\n")
	}
	return diags
}
