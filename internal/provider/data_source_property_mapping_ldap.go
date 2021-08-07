package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPPropertyMapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLDAPPropertyMappingRead,
		Description: "Get LDAP Property mappings",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"object_field": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLDAPPropertyMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.PropertymappingsApi.PropertymappingsLdapList(ctx)
	if n, ok := d.GetOk("name"); ok {
		req = req.Name(n.(string))
	}
	if m, ok := d.GetOk("managed"); ok {
		req = req.Managed(m.(string))
	}
	if m, ok := d.GetOk("object_field"); ok {
		req = req.ObjectField(m.(string))
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
	d.Set("object_field", f.ObjectField)
	return diags
}
