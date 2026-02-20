package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourcePropertyMappingProviderSCIM() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropertyMappingProviderSCIMRead,
		Description: "Customization --- Get SCIM Provider Property mappings",
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

func dataSourcePropertyMappingProviderSCIMRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PropertymappingsApi.PropertymappingsProviderScimList(ctx)

	if _, ok := d.GetOk("managed_list"); ok {
		req = req.Managed(helpers.CastSlice[string](d, "managed_list"))
	} else if m, ok := d.GetOk("managed"); ok {
		req = req.Managed([]string{m.(string)})
	}

	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
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
		helpers.SetWrapper(d, "ids", ids)
	} else {
		f := res.Results[0]
		d.SetId(f.Pk)
		helpers.SetWrapper(d, "name", f.Name)
		helpers.SetWrapper(d, "name", f.Name)
		helpers.SetWrapper(d, "expression", f.Expression)
	}
	return diags
}
