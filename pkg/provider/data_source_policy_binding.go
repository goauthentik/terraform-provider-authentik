package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func dataSourcePolicyBinding() *schema.Resource {
	return &schema.Resource{
		Description: "Customization --- Get policy bindings by id or target",
		ReadContext: dataSourcePolicyBindingRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"id", "target"},
			},
			"target": {
				Type:         schema.TypeString,
				Description:  "ID of the object this binding should apply to",
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"id", "target"},
			},
			"group": {
				Type:          schema.TypeString,
				Description:   "UUID of the group",
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"policy", "user"},
				RequiredWith:  []string{"target"},
			},
			"policy": {
				Type:          schema.TypeString,
				Description:   "UUID of the policy",
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"group", "user"},
				RequiredWith:  []string{"target"},
			},
			"user": {
				Type:          schema.TypeInt,
				Description:   "PK of the user",
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"group", "policy"},
				RequiredWith:  []string{"target"},
			},
			"order": {
				Type:         schema.TypeInt,
				Description:  "Order of the policy binding within the target",
				Computed:     true,
				Optional:     true,
				RequiredWith: []string{"target"},
			},

			// General attributes
			"negate": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failure_result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyBindingRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, idOk := d.GetOk("id")

	target, targetOk := d.GetOk("target")
	group, groupOk := d.GetOk("group")
	policy, policyOk := d.GetOk("policy")
	user, userOk := d.GetOk("user")
	order, orderOk := d.GetOk("order")

	if !idOk && (!targetOk || (!groupOk && !policyOk && !userOk)) {
		return diag.Errorf("Neither id nor target and user/group/policy were provided")
	}

	var pbs []api.PolicyBinding

	if idOk {
		req := c.client.PoliciesAPI.PoliciesBindingsRetrieve(ctx, id.(string))

		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		pbs = append(pbs, *res)
	} else {
		req := c.client.PoliciesAPI.PoliciesBindingsList(ctx)

		req = req.Target(target.(string))

		res, hr, err := req.Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}

		for _, v := range res.Results {
			if groupOk {
				if v.Group.IsSet() && v.GetGroup() == group.(string) {
					pbs = append(pbs, v)
				}
			}
			if policyOk {
				if v.Policy.IsSet() && v.GetPolicy() == policy.(string) {
					pbs = append(pbs, v)
				}
			}
			if userOk {
				if v.User.IsSet() && int(v.GetUser()) == user.(int) {
					pbs = append(pbs, v)
				}
			}
		}
	}

	if orderOk {
		n := 0
		for _, v := range pbs {
			if int(v.GetOrder()) == order.(int) {
				pbs[n] = v
				n++
			}
		}
		pbs = pbs[:n]
	}

	if len(pbs) < 1 {
		return diag.Errorf("No matching policy bindings found")
	}

	if len(pbs) > 1 {
		return diag.Errorf("Multiple matching policy bindings found. Use order to select one.")
	}

	f := pbs[0]
	d.SetId(f.Pk)
	helpers.SetWrapper(d, "target", f.Target)
	helpers.SetWrapper(d, "policy", f.Policy.Get())
	helpers.SetWrapper(d, "user", f.User.Get())
	helpers.SetWrapper(d, "group", f.Group.Get())
	helpers.SetWrapper(d, "order", f.Order)
	helpers.SetWrapper(d, "negate", f.Negate)
	helpers.SetWrapper(d, "enabled", f.Enabled)
	helpers.SetWrapper(d, "timeout", f.Timeout)
	helpers.SetWrapper(d, "failure_result", f.FailureResult)
	return diags
}
