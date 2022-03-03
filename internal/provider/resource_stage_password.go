package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func resourceStagePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStagePasswordCreate,
		ReadContext:   resourceStagePasswordRead,
		UpdateContext: resourceStagePasswordUpdate,
		DeleteContext: resourceStagePasswordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backends": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"configure_flow": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failed_attempts_before_cancel": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
		},
	}
}

func resourceStagePasswordSchemaToProvider(d *schema.ResourceData) *api.PasswordStageRequest {
	r := api.PasswordStageRequest{
		Name: d.Get("name").(string),
	}

	if s, sok := d.GetOk("configure_flow"); sok && s.(string) != "" {
		r.ConfigureFlow.Set(stringToPointer(s.(string)))
	}

	if fa, sok := d.GetOk("failed_attempts_before_cancel"); sok {
		r.FailedAttemptsBeforeCancel = intToPointer(fa.(int))
	}

	backend := make([]api.BackendsEnum, 0)
	for _, backendS := range d.Get("backends").([]interface{}) {
		backend = append(backend, api.BackendsEnum(backendS.(string)))
	}
	r.Backends = backend

	return &r
}

func resourceStagePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceStagePasswordSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPasswordCreate(ctx).PasswordStageRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePasswordRead(ctx, d, m)
}

func resourceStagePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.StagesApi.StagesPasswordRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.Set("name", res.Name)
	d.Set("backends", res.Backends)
	if res.ConfigureFlow.IsSet() {
		d.Set("configure_flow", res.ConfigureFlow.Get())
	}
	d.Set("failed_attempts_before_cancel", res.FailedAttemptsBeforeCancel)
	return diags
}

func resourceStagePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourceStagePasswordSchemaToProvider(d)

	res, hr, err := c.client.StagesApi.StagesPasswordUpdate(ctx, d.Id()).PasswordStageRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceStagePasswordRead(ctx, d, m)
}

func resourceStagePasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.StagesApi.StagesPasswordDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
