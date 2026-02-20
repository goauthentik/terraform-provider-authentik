package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourcePolicyBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Customization --- ",
		CreateContext: resourcePolicyBindingCreate,
		ReadContext:   resourcePolicyBindingRead,
		UpdateContext: resourcePolicyBindingUpdate,
		DeleteContext: resourcePolicyBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the object this binding should apply to",
			},
			"policy": {
				Type:        schema.TypeString,
				Description: "UUID of the policy",
				Optional:    true,
			},
			"user": {
				Type:        schema.TypeInt,
				Description: "PK of the user",
				Optional:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Description: "UUID of the group",
				Optional:    true,
			},

			// General attributes
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"negate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"failure_result": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourcePolicyBindingSchemaToModel(d *schema.ResourceData) *api.PolicyBindingRequest {
	m := api.PolicyBindingRequest{
		Target:        d.Get("target").(string),
		Order:         int32(d.Get("order").(int)),
		Negate:        new(d.Get("negate").(bool)),
		Enabled:       new(d.Get("enabled").(bool)),
		Timeout:       new(int32(d.Get("timeout").(int))),
		FailureResult: new(d.Get("failure_result").(bool)),
		Policy:        *api.NewNullableString(helpers.GetP[string](d, "policy")),
		User:          *api.NewNullableInt32(helpers.GetIntP(d, ("user"))),
		Group:         *api.NewNullableString(helpers.GetP[string](d, "group")),
	}
	return &m
}

func resourcePolicyBindingCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyBindingSchemaToModel(d)

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsCreate(ctx).PolicyBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyBindingRead(ctx, d, m)
}

func resourcePolicyBindingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "target", res.Target)
	helpers.SetWrapper(d, "policy", res.Policy.Get())
	helpers.SetWrapper(d, "user", res.User.Get())
	helpers.SetWrapper(d, "group", res.Group.Get())
	helpers.SetWrapper(d, "order", res.Order)
	helpers.SetWrapper(d, "negate", res.Negate)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "timeout", res.Timeout)
	helpers.SetWrapper(d, "failure_result", res.FailureResult)
	return diags
}

func resourcePolicyBindingUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app := resourcePolicyBindingSchemaToModel(d)

	res, hr, err := c.client.PoliciesApi.PoliciesBindingsUpdate(ctx, d.Id()).PolicyBindingRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourcePolicyBindingRead(ctx, d, m)
}

func resourcePolicyBindingDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.PoliciesApi.PoliciesBindingsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
