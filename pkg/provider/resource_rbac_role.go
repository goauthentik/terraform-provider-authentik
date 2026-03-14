package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRBACRole() *schema.Resource {
	return &schema.Resource{
		Description:   "RBAC --- ",
		CreateContext: resourceRBACRoleCreate,
		ReadContext:   resourceRBACRoleRead,
		UpdateContext: resourceRBACRoleUpdate,
		DeleteContext: resourceRBACRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRBACRoleSchemaToModel(d *schema.ResourceData) (*api.RoleRequest, diag.Diagnostics) {
	m := api.RoleRequest{
		Name: d.Get("name").(string),
	}
	return &m, nil
}

func resourceRBACRoleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*helpers.APIClient)

	app, diags := resourceRBACRoleSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.Client.RbacApi.RbacRolesCreate(ctx).RoleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRBACRoleRead(ctx, d, m)
}

func resourceRBACRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*helpers.APIClient)

	res, hr, err := c.Client.RbacApi.RbacRolesRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	return diags
}

func resourceRBACRoleUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*helpers.APIClient)

	app, di := resourceRBACRoleSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.Client.RbacApi.RbacRolesUpdate(ctx, d.Id()).RoleRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceRBACRoleRead(ctx, d, m)
}

func resourceRBACRoleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*helpers.APIClient)
	hr, err := c.Client.RbacApi.RbacRolesDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
