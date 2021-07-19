package authentik

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFlow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowRead,
		Description: "Get flows by Slug and/or designation",
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
		},
	}
}

func dataSourceFlowRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)

	req := c.client.FlowsApi.FlowsInstancesList(ctx)
	if s, ok := d.GetOk("slug"); ok {
		req.Slug(s.(string))
	}
	if d, ok := d.GetOk("designation"); ok {
		req.Designation(d.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching flows found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	d.Set("title", f.Title)
	d.Set("name", f.Name)
	d.Set("slug", f.Slug)
	d.Set("designation", f.Designation)
	return diags
}
