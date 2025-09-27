package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
)

func resourceRBACInitialPermissions() *schema.Resource {
	return &schema.Resource{
		Description:   "RBAC --- ",
		CreateContext: resourceRBACInitialPermissionsCreate,
		ReadContext:   resourceRBACInitialPermissionsRead,
		UpdateContext: resourceRBACInitialPermissionsUpdate,
		DeleteContext: resourceRBACInitialPermissionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      helpers.EnumToDescription(api.AllowedInitialPermissionsModeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedInitialPermissionsModeEnumEnumValues),
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceRBACInitialPermissionsSchemaToModel(d *schema.ResourceData) (*api.InitialPermissionsRequest, diag.Diagnostics) {
	m := api.InitialPermissionsRequest{
		Name:        d.Get("name").(string),
		Role:        d.Get("role").(string),
		Mode:        api.InitialPermissionsModeEnum(d.Get("mode").(string)),
		Permissions: castSliceInt32(d.Get("permissions").([]interface{})),
	}
	return &m, nil
}

func resourceRBACInitialPermissionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceRBACInitialPermissionsSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.RbacApi.RbacInitialPermissionsCreate(ctx).InitialPermissionsRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceRBACInitialPermissionsRead(ctx, d, m)
}

func resourceRBACInitialPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.RbacApi.RbacInitialPermissionsRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "role", res.Role)
	setWrapper(d, "mode", res.Mode)
	localPermissions := castSlice[int](d.Get("permissions").([]interface{}))
	setWrapper(d, "permissions", helpers.ListConsistentMerge(localPermissions, slice32ToInt(res.Permissions)))
	return diags
}

func resourceRBACInitialPermissionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	app, di := resourceRBACInitialPermissionsSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.RbacApi.RbacInitialPermissionsUpdate(ctx, int32(id)).InitialPermissionsRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceRBACInitialPermissionsRead(ctx, d, m)
}

func resourceRBACInitialPermissionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.RbacApi.RbacInitialPermissionsDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
