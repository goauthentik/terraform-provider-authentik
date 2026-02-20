package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceOutpostRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	if id, ok := d.GetOk("id"); ok {
		res, hr, err := c.client.OutpostsApi.OutpostsInstancesRetrieve(ctx, id.(string)).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		d.SetId(res.Pk)
		helpers.SetWrapper(d, "name", res.Name)
		return nil
	}

	if name, ok := d.GetOk("name"); ok {
		res, hr, err := c.client.OutpostsApi.OutpostsInstancesList(ctx).NameIexact(name.(string)).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		if len(res.Results) < 1 {
			return diag.Errorf("No matching outpost found")
		}
		if len(res.Results) > 1 {
			return diag.Errorf("Multiple outposts found")
		}
		d.SetId(res.Results[0].Pk)
		helpers.SetWrapper(d, "name", res.Results[0].Name)
		return nil
	}

	return diag.Errorf("Neither id nor name were provided")
}
