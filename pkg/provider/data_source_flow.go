package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceFlow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowRead,
		Description: "Flows & Stages --- Get flows by Slug and/or designation",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"designation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"authentication": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceFlowRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.FlowsApi.FlowsInstancesList(ctx)
	if s, ok := d.GetOk("slug"); ok {
		req = req.Slug(s.(string))
	}
	if des, ok := d.GetOk("designation"); ok {
		req = req.Designation(des.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching flows found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	helpers.SetWrapper(d, "title", f.Title)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetWrapper(d, "slug", f.Slug)
	helpers.SetWrapper(d, "designation", f.Designation)
	helpers.SetWrapper(d, "authentication", f.Authentication)
	return diags
}
