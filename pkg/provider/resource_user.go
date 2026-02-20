package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"type": {
				Type:             schema.TypeString,
				Default:          api.USERTYPEENUM_INTERNAL,
				Optional:         true,
				Description:      helpers.EnumToDescription(api.AllowedUserTypeEnumEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedUserTypeEnumEnumValues),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Optionally set the user's password. Changing the password in authentik will not trigger an update here.`,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Default:  "users",
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
		},
	}
}

func resourceUserSchemaToModel(d *schema.ResourceData) (*api.UserRequest, diag.Diagnostics) {
	m := api.UserRequest{
		Name:     d.Get("name").(string),
		Username: d.Get("username").(string),
		Type:     api.UserTypeEnum(d.Get("type").(string)).Ptr(),
		IsActive: new(d.Get("is_active").(bool)),
		Path:     new(d.Get("path").(string)),
		Email:    helpers.GetP[string](d, "email"),
		Groups:   helpers.CastSlice[string](d, "groups"),
		Roles:    helpers.CastSlice[string](d, "roles"),
	}
	attr, err := helpers.GetJSON[map[string]any](d, ("attributes"))
	m.Attributes = attr
	return &m, err
}

func resourceUserSetPassword(d *schema.ResourceData, c *APIClient, ctx context.Context) diag.Diagnostics {
	password, ok := d.Get("password").(string)
	if !ok || password == "" {
		return nil
	}
	if !d.IsNewResource() {
		return nil
	}
	uid, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.CoreApi.CoreUsersSetPasswordCreate(ctx, int32(uid)).UserPasswordSetRequest(api.UserPasswordSetRequest{
		Password: password,
	}).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	helpers.SetWrapper(d, "password", password)
	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceUserSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreUsersCreate(ctx).UserRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))

	diags = resourceUserSetPassword(d, c, ctx)
	if diags != nil {
		return diags
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.CoreApi.CoreUsersRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "type", res.Type)
	helpers.SetWrapper(d, "username", res.Username)
	helpers.SetWrapper(d, "email", res.Email)
	helpers.SetWrapper(d, "is_active", res.IsActive)
	helpers.SetWrapper(d, "path", res.Path)
	helpers.SetWrapper(d, "groups", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "groups"),
		res.Groups,
	))
	helpers.SetWrapper(d, "roles", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "roles"),
		res.Roles,
	))
	return helpers.SetJSON(d, "attributes", res.Attributes)
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)

	app, di := resourceUserSchemaToModel(d)
	if di != nil {
		return di
	}
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.CoreApi.CoreUsersUpdate(ctx, int32(id)).UserRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))

	diags := resourceUserSetPassword(d, c, ctx)
	if diags != nil {
		return diags
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.CoreApi.CoreUsersDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
