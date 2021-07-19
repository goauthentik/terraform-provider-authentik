package authentik

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Description: "Get groups by name",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)

	req := c.client.CoreApi.CoreGroupsList(ctx)
	if n, ok := d.GetOk("name"); ok {
		req.Name(n.(string))
	}
	if i, ok := d.GetOk("is_superuser"); ok {
		req.IsSuperuser(i.(bool))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching groups found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	d.Set("name", f.Name)
	d.Set("is_superuser", f.IsSuperuser)
	return diags
}
