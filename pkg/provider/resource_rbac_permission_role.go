package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		UpdateContext: resourceRBACRoleObjectPermissionUpdate,
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
		r.ObjectPk = api.PtrString(d.Get("object_id").(string))
	}
	return &r
}

func resourceRBACRoleObjectPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceRBACRoleObjectPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	_, modelSet := d.GetOk("model")
	if !modelSet {
		splitCodename := strings.Split(d.Get("permission").(string), ".")
		res, hr, err := c.client.RbacApi.RbacPermissionsList(ctx).Codename(splitCodename[1]).Role(d.Get("role").(string)).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		if len(res.Results) < 1 {
			return diag.Errorf("Permission not found")
		}
		helpers.SetWrapper(d, "permission", fmt.Sprintf("%s.%s", res.Results[0].AppLabel, res.Results[0].Codename))
	} else {
		res, hr, err := c.client.RbacApi.RbacPermissionsRolesRetrieve(ctx, int32(id)).Execute()
		if err != nil {
			return helpers.HTTPToDiag(d, hr, err)
		}
		helpers.SetWrapper(d, "object_id", res.ObjectPk)
	}
	return diags
}

func resourceRBACRoleObjectPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	req := resourceRBACRoleObjectPermissionSchemaToProvider(d).ObjectPk
	res, hr, err := c.client.RbacApi.RbacPermissionsRolesUpdate(ctx, int32(id)).ExtraRoleObjectPermissionRequest(api.ExtraRoleObjectPermissionRequest{
		ObjectPk: *req,
	}).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Id)))
	return resourceRBACRoleObjectPermissionRead(ctx, d, m)
}

func resourceRBACRoleObjectPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.RbacApi.RbacPermissionsRolesDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
