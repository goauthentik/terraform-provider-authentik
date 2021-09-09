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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"managed_list"},
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"managed_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Retrive multiple property mappings",
			},

			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of ids when `managed_list` is set.",
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

	if ml, ok := d.GetOk("managed_list"); ok {
		req = req.Managed(sliceToStringPointer(ml.([]interface{})))
	} else if m, ok := d.GetOk("managed"); ok {
		mm := m.(string)
		req = req.Managed([]*string{&mm})
	}

	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
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
	if _, ok := d.GetOk("managed_list"); ok {
		d.SetId("-1")
		ids := make([]string, len(res.Results))
		for i, r := range res.Results {
			ids[i] = r.Pk
		}
		d.Set("ids", ids)
	} else {
		f := res.Results[0]
		d.SetId(f.Pk)
		d.Set("name", f.Name)
		d.Set("expression", f.Expression)
		d.Set("scope_name", f.ScopeName)
		d.Set("description", f.Description)
	}
	return diags
}
