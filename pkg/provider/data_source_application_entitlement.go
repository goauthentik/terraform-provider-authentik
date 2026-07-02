package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceApplicationEntitlement() *schema.Resource {
	return &schema.Resource{
		Description: "Applications --- Get application entitlements by id or application uuid and entitlement name",
		ReadContext: dataSourceApplicationEntitlementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"app", "name"},
			},
			"app": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"app"},
			},
			"name": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"id"},
				RequiredWith:  []string{"name"},
			},
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApplicationEntitlementRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, idOk := d.GetOk("id")
	app, appOk := d.GetOk("app")
	name, nameOk := d.GetOk("name")

	if !idOk && (!appOk || !nameOk) {
		return diag.Errorf("Neither id nor app/name pair were provided")
	}

	req := c.client.CoreAPI.CoreApplicationEntitlementsList(ctx)

	if idOk {
		req = req.PbmUuid(id.(string))
	} else {
		req = req.App(app.(string))
		req = req.Name(name.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No matching application entitlements found")
	}

	if len(res.Results) > 1 {
		// In theory, impossible..
		return diag.Errorf("Multiple application entitlements found")
	}

	f := res.Results[0]
	d.SetId(f.PbmUuid)
	helpers.SetWrapper(d, "app", f.App)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetJSON(d, "attributes", f.Attributes)
	return diags
}
