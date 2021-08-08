package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCertificateKeyPair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCertificateKeyPairRead,
		Description: "Get certificate-key pairs by name",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"key_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceCertificateKeyPairRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.CryptoApi.CryptoCertificatekeypairsList(ctx)
	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	if len(res.Results) < 1 {
		return diag.Errorf("No matching groups found")
	}
	f := res.Results[0]

	d.SetId(f.Pk)
	d.Set("name", f.Name)
	d.Set("expiry", f.CertExpiry.String())
	d.Set("subject", f.CertSubject)
	d.Set("fingerprint1", f.FingerprintSha1)
	d.Set("fingerprint256", f.FingerprintSha256)

	rc, hr, err := c.client.CryptoApi.CryptoCertificatekeypairsViewCertificateRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	d.Set("certificate_data", rc.Data+"\n")

	rk, _, err := c.client.CryptoApi.CryptoCertificatekeypairsViewPrivateKeyRetrieve(ctx, d.Id()).Execute()
	if err == nil {
		d.Set("key_data", rk.Data+"\n")
	}
	return diags
}
