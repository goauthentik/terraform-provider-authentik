package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourceProviderSAMLMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProviderSAMLMetadataRead,
		Description: "Applications --- Get SAML Provider metadata",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"provider_id"},
				Description:   "Find provider by name",
			},
			"provider_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "Find provider by ID",
			},

			"metadata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SAML Metadata",
			},
		},
	}
}

func dataSourceProviderSAMLMetadataRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, ok := d.GetOk("provider_id")
	if !ok {
		req := c.client.ProvidersApi.ProvidersSamlList(ctx)
		if m, ok := d.Get("name").(string); ok {
			req = req.Name(m)
		}
		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		if len(res.Results) < 1 {
			return diag.Errorf("no matching providers found")
		}
		id = int(res.Results[0].Pk)
	}
	finalId := int32(id.(int))
	d.SetId(strconv.FormatInt(int64(finalId), 10))

	meta, hr, err := c.client.ProvidersApi.ProvidersSamlMetadataRetrieve(ctx, finalId).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "metadata", meta.Metadata)
	return diags
}
