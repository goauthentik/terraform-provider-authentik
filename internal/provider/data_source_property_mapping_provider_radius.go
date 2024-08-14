package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyMappingProviderRadius() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropertyMappingProviderRadiusRead,
		Description: "Customization --- Get Radius Provider Property mappings",
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
				Description: "Retrieve multiple property mappings",
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

			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePropertyMappingProviderRadiusRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PropertymappingsApi.PropertymappingsProviderRadiusList(ctx)

	if ml, ok := d.GetOk("managed_list"); ok {
		req = req.Managed(castSlice[string](ml.([]interface{})))
	} else if m, ok := d.GetOk("managed"); ok {
		req = req.Managed([]string{m.(string)})
	}

	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
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
		setWrapper(d, "ids", ids)
	} else {
		f := res.Results[0]
		d.SetId(f.Pk)
		setWrapper(d, "name", f.Name)
		setWrapper(d, "name", f.Name)
		setWrapper(d, "expression", f.Expression)
	}
	return diags
}
