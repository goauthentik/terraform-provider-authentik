package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceRBACRoleObjectPermission() *schema.Resource {
	return &schema.Resource{
		Description:   "RBAC --- ",
		CreateContext: resourceRBACRoleObjectPermissionCreate,
		ReadContext:   resourceRBACRoleObjectPermissionRead,
		// UpdateContext: resourceRBACRoleObjectPermissionUpdate,
		DeleteContext: resourceRBACRoleObjectPermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Description:      helpers.EnumToDescription(api.AllowedModelEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedModelEnumEnumValues),
			},
			"object_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRBACRoleObjectPermissionSchemaToProvider(d *schema.ResourceData) *api.PermissionAssignRequest {
	r := api.PermissionAssignRequest{
		Permissions: []string{d.Get("permission").(string)},
	}
	if d.Get("model").(string) != "" {
		r.Model = api.ModelEnum(d.Get("model").(string)).Ptr()
	}
	if d.Get("object_id").(string) != "" {
		r.ObjectPk = new(d.Get("object_id").(string))
	}
	return &r
}

func resourceRBACRoleObjectPermissionCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceRBACRoleObjectPermissionSchemaToProvider(d)

	role := d.Get("role").(string)

	res, hr, err := c.client.RbacApi.RbacPermissionsAssignedByRolesAssign(ctx, role).PermissionAssignRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	if len(res) != 1 {
		return diag.Errorf("invalid API response")
	}
	d.SetId(res[0].Id)
	return resourceRBACRoleObjectPermissionRead(ctx, d, m)
}

func resourceRBACRoleObjectPermissionRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	_, object := d.GetOk("object_id")
	if object {
		perms, hr, err := helpers.Paginator(c.client.RbacApi.RbacPermissionsRolesList(ctx).Uuid(d.Get("role").(string)), helpers.PaginatorOptions{})
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		for _, perm := range perms {
			if perm.Id == int32(id) {
				helpers.SetWrapper(d, "permission", fmt.Sprintf("%s.%s", perm.AppLabel, perm.Codename))
				helpers.SetWrapper(d, "object_id", perm.ObjectPk)
				return diags
			}
		}
	} else {
		perms, hr, err := helpers.Paginator(c.client.RbacApi.RbacPermissionsList(ctx).Role(d.Get("role").(string)), helpers.PaginatorOptions{})
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		for _, perm := range perms {
			fqpn := fmt.Sprintf("%s.%s", perm.AppLabel, perm.Codename)
			if fqpn != d.Get("permission").(string) {
				continue
			}
			helpers.SetWrapper(d, "permission", fqpn)
			return diags
		}
	}
	return diag.FromErr(errors.New("permission not found"))
}

func resourceRBACRoleObjectPermissionDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	req := api.PatchedPermissionAssignRequest{
		Permissions: []string{d.Get("permission").(string)},
	}
	if d.Get("model").(string) != "" {
		req.Model = api.ModelEnum(d.Get("model").(string)).Ptr()
	}
	if d.Get("object_id").(string) != "" {
		req.ObjectPk = new(d.Get("object_id").(string))
	}

	hr, err := c.client.RbacApi.RbacPermissionsAssignedByRolesUnassignPartialUpdate(ctx, d.Get("role").(string)).PatchedPermissionAssignRequest(req).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
