package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataOutpostServiceConnectionsKubernetes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataOutpostServiceConnectionsKubernetesRead,
		Description: "Applications --- Get a Kubernetes Service Connection by name",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Optional:  true,
				Sensitive: true,
			},
			"local": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataOutpostServiceConnectionsKubernetesRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	req := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesList(ctx)
	if s, ok := d.GetOk("name"); ok {
		req = req.Name(s.(string))
	}

	res, hr, err := req.Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	if len(res.Results) < 1 {
		return diag.Errorf("No Kubernetes Outpost Service Connections found")
	}
	f := res.Results[0]
	d.SetId(f.Pk)
	helpers.SetWrapper(d, "name", f.Name)
	helpers.SetWrapper(d, "local", f.Local)
	helpers.SetWrapper(d, "verify_ssl", f.VerifySsl)
	helpers.SetJSON(d, "kubeconfig", f.Kubeconfig)
	return diags
}
