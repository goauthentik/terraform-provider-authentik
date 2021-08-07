package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceScopeMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScopeMappingRead,
		Description: "Get OAuth Scope mappings",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceScopeMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PropertymappingsApi.PropertymappingsScopeList(ctx)
	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}
	if m, ok := d.GetOk("managed"); ok {
		req = req.Managed(m.(string))
	}
	if m, ok := d.GetOk("scope_name"); ok {
		req = req.ScopeName(m.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching mappings found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	d.Set("name", f.Name)
	d.Set("expression", f.Expression)
	d.Set("scope_name", f.ScopeName)
	d.Set("description", f.Description)
	return diags
}
