package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourcePolicyEventMatcher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyEventMatcherCreate,
		ReadContext:   resourcePolicyEventMatcherRead,
		UpdateContext: resourcePolicyEventMatcherUpdate,
		DeleteContext: resourcePolicyEventMatcherDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"action": {
				Type: schema.TypeString,
				// TODO: Fix schema not allowing blank values
				Required: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app": {
				Type: schema.TypeString,
				// TODO: Fix schema not allowing blank values
				Required: true,
			},
		},
	}
}

func resourcePolicyEventMatcherSchemaToProvider(d *schema.ResourceData) *api.EventMatcherPolicyRequest {
	r := api.EventMatcherPolicyRequest{
		Name:             d.Get("name").(string),
		ExecutionLogging: boolToPointer(d.Get("execution_logging").(bool)),
	}

	act := api.EventActions(d.Get("action").(string))
	r.Action.Set(&act)

	if p, pSet := d.GetOk("client_ip"); pSet {
		r.ClientIp = stringToPointer(p.(string))
	}

	app := api.AppEnum(d.Get("app").(string))
	r.App.Set(&app)

	return &r
}

func resourcePolicyEventMatcherCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourcePolicyEventMatcherSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherCreate(ctx).EventMatcherPolicyRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyEventMatcherRead(ctx, d, m)
}

func resourcePolicyEventMatcherRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "execution_logging", res.ExecutionLogging)
	setWrapper(d, "action", res.Action.Get())
	setWrapper(d, "client_ip", res.ClientIp)
	setWrapper(d, "app", res.App.Get())
	return diags
}

func resourcePolicyEventMatcherUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyEventMatcherSchemaToProvider(d)

	res, hr, err := c.client.PoliciesApi.PoliciesEventMatcherUpdate(ctx, d.Id()).EventMatcherPolicyRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyEventMatcherRead(ctx, d, m)
}

func resourcePolicyEventMatcherDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesEventMatcherDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
