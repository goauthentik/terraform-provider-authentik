package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceServiceConnectionKubernetes() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceServiceConnectionKubernetesCreate,
		ReadContext:   resourceServiceConnectionKubernetesRead,
		UpdateContext: resourceServiceConnectionKubernetesUpdate,
		DeleteContext: resourceServiceConnectionKubernetesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"local": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"kubeconfig": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceServiceConnectionKubernetesSchemaToModel(d *schema.ResourceData) (*api.KubernetesServiceConnectionRequest, diag.Diagnostics) {
	m := api.KubernetesServiceConnectionRequest{
		Name:      d.Get("name").(string),
		VerifySsl: new(d.Get("verify_ssl").(bool)),
		Local:     new(d.Get("local").(bool)),
	}

	attr, err := helpers.GetJSON[map[string]any](d, ("kubeconfig"))
	m.Kubeconfig = attr
	return &m, err
}

func resourceServiceConnectionKubernetesCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceServiceConnectionKubernetesSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesCreate(ctx).KubernetesServiceConnectionRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionKubernetesRead(ctx, d, m)
}

func resourceServiceConnectionKubernetesRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "local", res.Local)
	helpers.SetWrapper(d, "verify_ssl", res.VerifySsl)
	return helpers.SetJSON(d, "kubeconfig", res.Kubeconfig)
}

func resourceServiceConnectionKubernetesUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceServiceConnectionKubernetesSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesUpdate(ctx, d.Id()).KubernetesServiceConnectionRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionKubernetesRead(ctx, d, m)
}

func resourceServiceConnectionKubernetesDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
