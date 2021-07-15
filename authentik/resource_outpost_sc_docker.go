package authentik

import (
	"context"

	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceConnectionDocker() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceConnectionDockerCreate,
		ReadContext:   resourceServiceConnectionDockerRead,
		UpdateContext: resourceServiceConnectionDockerUpdate,
		DeleteContext: resourceServiceConnectionDockerDelete,
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
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tls_verification": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_authentication": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceServiceConnectionDockerSchemaToModel(d *schema.ResourceData) (*api.DockerServiceConnectionRequest, diag.Diagnostics) {
	m := api.DockerServiceConnectionRequest{
		Name: d.Get("name").(string),
		Url:  d.Get("url").(string),
	}

	local := d.Get("local").(bool)
	m.Local = &local

	if l, ok := d.Get("tls_verification").(string); ok {
		m.TlsVerification.Set(&l)
	} else {
		m.TlsVerification.Set(nil)
	}

	if l, ok := d.Get("tls_authentication").(string); ok {
		m.TlsAuthentication.Set(&l)
	} else {
		m.TlsAuthentication.Set(nil)
	}
	return &m, nil
}

func resourceServiceConnectionDockerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	app, diags := resourceServiceConnectionDockerSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerCreate(ctx).DockerServiceConnectionRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionDockerRead(ctx, d, m)
}

func resourceServiceConnectionDockerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*ProviderAPIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.Set("name", res.Name)
	d.Set("url", res.Url)
	d.Set("local", res.Local)
	d.Set("tls_verification", res.TlsVerification)
	d.Set("tls_authentication", res.TlsAuthentication)
	return diags
}

func resourceServiceConnectionDockerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)

	app, di := resourceServiceConnectionDockerSchemaToModel(d)
	if di != nil {
		return di
	}

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerUpdate(ctx, d.Id()).DockerServiceConnectionRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionDockerRead(ctx, d, m)
}

func resourceServiceConnectionDockerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*ProviderAPIClient)
	hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr)
	}
	return diag.Diagnostics{}
}
