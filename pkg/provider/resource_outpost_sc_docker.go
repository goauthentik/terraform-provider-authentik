package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceServiceConnectionDocker() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
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
				Optional: true,
				Default:  "http+unix:///var/run/docker.sock",
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

func resourceServiceConnectionDockerSchemaToModel(d *schema.ResourceData) *api.DockerServiceConnectionRequest {
	m := api.DockerServiceConnectionRequest{
		Name:              d.Get("name").(string),
		Url:               d.Get("url").(string),
		Local:             new(d.Get("local").(bool)),
		TlsVerification:   *api.NewNullableString(helpers.GetP[string](d, "tls_verification")),
		TlsAuthentication: *api.NewNullableString(helpers.GetP[string](d, "tls_authentication")),
	}
	return &m
}

func resourceServiceConnectionDockerCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceServiceConnectionDockerSchemaToModel(d)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerCreate(ctx).DockerServiceConnectionRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionDockerRead(ctx, d, m)
}

func resourceServiceConnectionDockerRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "url", res.Url)
	helpers.SetWrapper(d, "local", res.Local)
	helpers.SetWrapper(d, "tls_verification", res.TlsVerification.Get())
	helpers.SetWrapper(d, "tls_authentication", res.TlsAuthentication.Get())
	return diags
}

func resourceServiceConnectionDockerUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceServiceConnectionDockerSchemaToModel(d)

	res, hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerUpdate(ctx, d.Id()).DockerServiceConnectionRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceServiceConnectionDockerRead(ctx, d, m)
}

func resourceServiceConnectionDockerDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.OutpostsApi.OutpostsServiceConnectionsDockerDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
