package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func resourceSourceLDAP() *schema.Resource {
	return &schema.Resource{
		Description:   "Directory --- ",
		CreateContext: resourceSourceLDAPCreate,
		ReadContext:   resourceSourceLDAPRead,
		UpdateContext: resourceSourceLDAPUpdate,
		DeleteContext: resourceSourceLDAPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: helpers.ModelSchema(
			api.MODELENUM_SOURCES_LDAP_LDAPSOURCE,
			map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"uuid": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"slug": {
					Type:     schema.TypeString,
					Required: true,
				},
				"user_path_template": {
					Type:     schema.TypeString,
					Default:  "goauthentik.io/sources/%(slug)s",
					Optional: true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},

				"server_uri": {
					Type:     schema.TypeString,
					Required: true,
				},
				"bind_cn": {
					Type:     schema.TypeString,
					Required: true,
				},
				"bind_password": {
					Type:      schema.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"start_tls": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sni": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"base_dn": {
					Type:     schema.TypeString,
					Required: true,
				},
				"additional_user_dn": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"additional_group_dn": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"user_object_filter": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "(objectClass=person)",
				},
				"user_membership_attribute": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "distinguishedName",
				},
				"group_object_filter": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "(objectClass=group)",
				},
				"group_membership_field": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "member",
				},
				"object_uniqueness_field": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "objectSid",
				},
				"lookup_groups_from_user": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sync_users": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sync_users_password": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sync_groups": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
				},
				"sync_parent_group": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password_login_update_internal_password": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"delete_not_found_objects": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"property_mappings": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"property_mappings_group": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		),
	}
}

func resourceSourceLDAPSchemaToSource(d *schema.ResourceData) *api.LDAPSourceRequest {
	r := api.LDAPSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          api.PtrBool(d.Get("enabled").(bool)),
		UserPathTemplate: api.PtrString(d.Get("user_path_template").(string)),

		BaseDn:       d.Get("base_dn").(string),
		ServerUri:    d.Get("server_uri").(string),
		BindCn:       api.PtrString(d.Get("bind_cn").(string)),
		BindPassword: api.PtrString(d.Get("bind_password").(string)),
		StartTls:     api.PtrBool(d.Get("start_tls").(bool)),
		Sni:          api.PtrBool(d.Get("sni").(bool)),

		AdditionalUserDn:        api.PtrString(d.Get("additional_user_dn").(string)),
		AdditionalGroupDn:       api.PtrString(d.Get("additional_group_dn").(string)),
		UserObjectFilter:        api.PtrString(d.Get("user_object_filter").(string)),
		UserMembershipAttribute: api.PtrString(d.Get("user_membership_attribute").(string)),
		GroupObjectFilter:       api.PtrString(d.Get("group_object_filter").(string)),
		GroupMembershipField:    api.PtrString(d.Get("group_membership_field").(string)),
		ObjectUniquenessField:   api.PtrString(d.Get("object_uniqueness_field").(string)),

		SyncUsers:                           api.PtrBool(d.Get("sync_users").(bool)),
		SyncUsersPassword:                   api.PtrBool(d.Get("sync_users_password").(bool)),
		SyncGroups:                          api.PtrBool(d.Get("sync_groups").(bool)),
		SyncParentGroup:                     *api.NewNullableString(helpers.GetP[string](d, "sync_parent_group")),
		PasswordLoginUpdateInternalPassword: api.PtrBool(d.Get("password_login_update_internal_password").(bool)),
		DeleteNotFoundObjects:               api.PtrBool(d.Get("delete_not_found_objects").(bool)),
		LookupGroupsFromUser:                api.PtrBool(d.Get("lookup_groups_from_user").(bool)),
		UserPropertyMappings:                helpers.CastSlice[string](d, "property_mappings"),
		GroupPropertyMappings:               helpers.CastSlice[string](d, "property_mappings_group"),
	}
	return &r
}

func resourceSourceLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapCreate(ctx).LDAPSourceRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "slug", res.Slug)
	helpers.SetWrapper(d, "uuid", res.Pk)
	helpers.SetWrapper(d, "enabled", res.Enabled)
	helpers.SetWrapper(d, "user_path_template", res.UserPathTemplate)

	helpers.SetWrapper(d, "base_dn", res.BaseDn)
	helpers.SetWrapper(d, "server_uri", res.ServerUri)
	helpers.SetWrapper(d, "bind_cn", res.BindCn)
	helpers.SetWrapper(d, "start_tls", res.StartTls)
	helpers.SetWrapper(d, "sni", res.Sni)
	helpers.SetWrapper(d, "base_dn", res.BaseDn)
	helpers.SetWrapper(d, "additional_user_dn", res.AdditionalUserDn)
	helpers.SetWrapper(d, "additional_group_dn", res.AdditionalGroupDn)
	helpers.SetWrapper(d, "user_object_filter", res.UserObjectFilter)
	helpers.SetWrapper(d, "user_membership_attribute", res.UserMembershipAttribute)
	helpers.SetWrapper(d, "group_object_filter", res.GroupObjectFilter)
	helpers.SetWrapper(d, "group_membership_field", res.GroupMembershipField)
	helpers.SetWrapper(d, "object_uniqueness_field", res.ObjectUniquenessField)
	helpers.SetWrapper(d, "lookup_groups_from_user", res.LookupGroupsFromUser)
	helpers.SetWrapper(d, "sync_users", res.SyncUsers)
	helpers.SetWrapper(d, "sync_users_password", res.SyncUsersPassword)
	helpers.SetWrapper(d, "sync_groups", res.SyncGroups)
	helpers.SetWrapper(d, "password_login_update_internal_password", res.PasswordLoginUpdateInternalPassword)
	helpers.SetWrapper(d, "delete_not_found_objects", res.DeleteNotFoundObjects)
	helpers.SetWrapper(d, "sync_parent_group", res.SyncParentGroup.Get())
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.UserPropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.GroupPropertyMappings,
	))
	return diags
}

func resourceSourceLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapUpdate(ctx, d.Id()).LDAPSourceRequest(*app).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
