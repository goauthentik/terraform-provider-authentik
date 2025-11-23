package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceOutpost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOutpostRead,
		Description: "Applications --- Get outposts by id or name",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourceOutpostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	if id, ok := d.GetOk("id"); ok {
		res, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, id.(string)).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		d.SetId(res.Pk)
		d.Set("name", res.Name)
		return nil
	}

	if name, ok := d.GetOk("name"); ok {
		res, hr, err := c.client.OutpostsApi.OutpostsInstancesList(ctx).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		var found []api.Outpost
		for _, o := range res.Results {
			if o.Name == name.(string) {
				found = append(found, o)
			}
		}
		if len(found) < 1 {
			return diag.Errorf("No matching outpost found")
		}
		if len(found) > 1 {
			return diag.Errorf("Multiple outposts found")
		}
		d.SetId(found[0].Pk)
		d.Set("id", found[0].Pk)
		return nil
	}

	return diag.Errorf("Neither id nor name were provided")
}
