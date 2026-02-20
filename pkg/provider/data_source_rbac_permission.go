package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceRBACPermission() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRBACPermissionRead,
		Description: "RBAC --- Get a permission by codename",
		Schema: map[string]*schema.Schema{
			"codename": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRBACPermissionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.RbacApi.RbacPermissionsList(ctx)
	if codename, ok := d.GetOk("codename"); ok {
		req = req.Codename(codename.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching flows found")
	}
	f := res.Results[0]
	d.SetId(strconv.Itoa(int(f.Id)))
	helpers.SetWrapper(d, "app", f.AppLabel)
	helpers.SetWrapper(d, "model", f.Model)
	return diags
}
