package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      JSONDescription,
				DiffSuppressFunc: diffSuppressJSON,
				ValidateDiagFunc: ValidateJSON,
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
		VerifySsl: api.PtrBool(d.Get("verify_ssl").(bool)),
		Local:     api.PtrBool(d.Get("local").(bool)),
	}

	attr, err := getJSON[map[string]interface{}](d.Get("kubeconfig"))
	m.Kubeconfig = attr
	return &m, err
}

func resourceServiceConnectionKubernetesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceServiceConnectionKubernetesSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesCreate(ctx).KubernetesServiceConnectionRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionKubernetesRead(ctx, d, m)
}

func resourceServiceConnectionKubernetesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "local", res.Local)
	setWrapper(d, "verify_ssl", res.VerifySsl)
	b, err := json.Marshal(res.Kubeconfig)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "kubeconfig", string(b))
	return diags
}

func resourceServiceConnectionKubernetesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceServiceConnectionKubernetesSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesUpdate(ctx, d.Id()).KubernetesServiceConnectionRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionKubernetesRead(ctx, d, m)
}

func resourceServiceConnectionKubernetesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsKubernetesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
