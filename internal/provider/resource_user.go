package provider

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
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
				Description:      EnumToDescription(api.AllowedUserTypeEnumEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedUserTypeEnumEnumValues),
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
			"attributes": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      "JSON format expected. Use jsonencode() to pass objects.",
				DiffSuppressFunc: diffSuppressJSON,
			},
		},
	}
}

func resourceUserSchemaToModel(d *schema.ResourceData) (*api.UserRequest, diag.Diagnostics) {
	m := api.UserRequest{
		Name:     d.Get("name").(string),
		Username: d.Get("username").(string),
		Type:     api.UserTypeEnum(d.Get("type").(string)).Ptr(),
		IsActive: api.PtrBool(d.Get("is_active").(bool)),
		Path:     api.PtrString(d.Get("path").(string)),
	}

	if l, ok := d.Get("email").(string); ok {
		m.Email = &l
	}

	m.Groups = castSlice[string](d.Get("groups").([]interface{}))

	attr := make(map[string]interface{})
	if l, ok := d.Get("attributes").(string); ok && l != "" {
		err := json.NewDecoder(strings.NewReader(l)).Decode(&attr)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}
	m.Attributes = attr
	return &m, nil
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
		return httpToDiag(d, hr, err)
	}
	setWrapper(d, "password", password)
	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	app, diags := resourceUserSchemaToModel(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.CoreApi.CoreUsersCreate(ctx).UserRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))

	diags = resourceUserSetPassword(d, c, ctx)
	if diags != nil {
		return diags
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	res, hr, err := c.client.CoreApi.CoreUsersRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "type", res.Type)
	setWrapper(d, "username", res.Username)
	setWrapper(d, "email", res.Email)
	setWrapper(d, "is_active", res.IsActive)
	setWrapper(d, "path", res.Path)
	b, err := json.Marshal(res.Attributes)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "attributes", string(b))
	localGroups := castSlice[string](d.Get("groups").([]interface{}))
	setWrapper(d, "groups", listConsistentMerge(localGroups, res.Groups))
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))

	diags := resourceUserSetPassword(d, c, ctx)
	if diags != nil {
		return diags
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.CoreApi.CoreUsersDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
