package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSourceRead,
		Description: "Directory --- Get Source by name, slug or managed",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.SourcesApi.SourcesAllList(ctx)
	if s, ok := d.GetOk("slug"); ok {
		req = req.Slug(s.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		req = req.Name(name.(string))
	}
	if managed, ok := d.GetOk("managed"); ok {
		req = req.Managed(managed.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching flows found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	setWrapper(d, "name", f.Name)
	setWrapper(d, "slug", f.Slug)
	setWrapper(d, "managed", f.Managed.Get())
	setWrapper(d, "uuid", f.Pk)
	return diags
}
