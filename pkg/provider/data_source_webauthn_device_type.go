package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceWebAuthnDeviceType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebAuthnDeviceTypeRead,
		Description: "Flows & Stages --- ",
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aaguid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWebAuthnDeviceTypeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.StagesApi.StagesAuthenticatorWebauthnDeviceTypesList(ctx)
	if s, ok := d.GetOk("description"); ok {
		req = req.Description(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching device types found")
	}
	f := res.Results[0]
	d.SetId(f.Aaguid)
	helpers.SetWrapper(d, "aaguid", f.Aaguid)
	helpers.SetWrapper(d, "description", f.Description)
	return diags
}
