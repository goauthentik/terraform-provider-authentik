package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_superuser": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"parents": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGroupSchemaToModel(d *schema.ResourceData) (*api.GroupRequest, diag.Diagnostics) {
	m := api.GroupRequest{
		Name:        d.Get("name").(string),
		IsSuperuser: new(d.Get("is_superuser").(bool)),
		Parents:     helpers.CastSlice[string](d, "parents"),
		Users:       helpers.CastSliceInt32(d.Get("users").([]any)),
		Roles:       helpers.CastSlice[string](d, "roles"),
	}
	attr, err := helpers.GetJSON[map[string]any](d, ("attributes"))
	m.Attributes = attr
	return &m, err
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceGroupSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreGroupsCreate(ctx).GroupRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	res, hr, err := c.client.CoreApi.CoreGroupsRetrieve(ctx, d.Id()).IncludeUsers(false).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "is_superuser", res.IsSuperuser)
	helpers.SetWrapper(d, "parents", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "parents"),
		res.Parents,
	))
	helpers.SetWrapper(d, "users", helpers.ListConsistentMerge(
		helpers.CastSlice[int](d, "users"),
		helpers.Slice32ToInt(res.Users),
	))
	helpers.SetWrapper(d, "roles", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "role"),
		res.Roles,
	))
	return helpers.SetJSON(d, "attributes", res.Attributes)
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceGroupSchemaToModel(d)
	if di != nil {
		return di
	}
	res, hr, err := c.client.CoreApi.CoreGroupsUpdate(ctx, d.Id()).GroupRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Pk)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.CoreApi.CoreGroupsDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
