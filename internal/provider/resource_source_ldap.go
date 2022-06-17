package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func resourceSourceLDAP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceLDAPCreate,
		ReadContext:   resourceSourceLDAPRead,
		UpdateContext: resourceSourceLDAPUpdate,
		DeleteContext: resourceSourceLDAPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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
	}
}

func resourceSourceLDAPSchemaToSource(d *schema.ResourceData) *api.LDAPSourceRequest {
	r := api.LDAPSourceRequest{
		Name:             d.Get("name").(string),
		Slug:             d.Get("slug").(string),
		Enabled:          boolToPointer(d.Get("enabled").(bool)),
		UserPathTemplate: stringToPointer(d.Get("user_path_template").(string)),

		BaseDn:       d.Get("base_dn").(string),
		ServerUri:    d.Get("server_uri").(string),
		BindCn:       stringToPointer(d.Get("bind_cn").(string)),
		BindPassword: stringToPointer(d.Get("bind_password").(string)),
		StartTls:     boolToPointer(d.Get("start_tls").(bool)),

		AdditionalUserDn:      stringToPointer(d.Get("additional_user_dn").(string)),
		AdditionalGroupDn:     stringToPointer(d.Get("additional_group_dn").(string)),
		UserObjectFilter:      stringToPointer(d.Get("user_object_filter").(string)),
		GroupObjectFilter:     stringToPointer(d.Get("group_object_filter").(string)),
		GroupMembershipField:  stringToPointer(d.Get("group_membership_field").(string)),
		ObjectUniquenessField: stringToPointer(d.Get("object_uniqueness_field").(string)),

		SyncUsers:         boolToPointer(d.Get("sync_users").(bool)),
		SyncUsersPassword: boolToPointer(d.Get("sync_users_password").(bool)),
		SyncGroups:        boolToPointer(d.Get("sync_groups").(bool)),
	}

	if s, sok := d.GetOk("sync_parent_group"); sok && s.(string) != "" {
		r.SyncParentGroup.Set(stringToPointer(s.(string)))
	}

	r.PropertyMappings = sliceToString(d.Get("property_mappings").([]interface{}))
	r.PropertyMappingsGroup = sliceToString(d.Get("property_mappings_group").([]interface{}))
	return &r
}

func resourceSourceLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapCreate(ctx).LDAPSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "slug", res.Slug)
	setWrapper(d, "uuid", res.Pk)
	setWrapper(d, "enabled", res.Enabled)
	setWrapper(d, "user_path_template", res.UserPathTemplate)

	setWrapper(d, "base_dn", res.BaseDn)
	setWrapper(d, "server_uri", res.ServerUri)
	setWrapper(d, "bind_cn", res.BindCn)
	setWrapper(d, "start_tls", res.StartTls)
	setWrapper(d, "base_dn", res.BaseDn)
	setWrapper(d, "additional_user_dn", res.AdditionalUserDn)
	setWrapper(d, "additional_group_dn", res.AdditionalGroupDn)
	setWrapper(d, "user_object_filter", res.UserObjectFilter)
	setWrapper(d, "group_object_filter", res.GroupObjectFilter)
	setWrapper(d, "group_membership_field", res.GroupMembershipField)
	setWrapper(d, "object_uniqueness_field", res.ObjectUniquenessField)
	setWrapper(d, "sync_users", res.SyncUsers)
	setWrapper(d, "sync_users_password", res.SyncUsersPassword)
	setWrapper(d, "sync_groups", res.SyncGroups)
	if res.SyncParentGroup.IsSet() {
		setWrapper(d, "sync_parent_group", res.SyncParentGroup.Get())
	}
	localMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	setWrapper(d, "property_mappings", stringListConsistentMerge(localMappings, res.PropertyMappings))
	localGroupMappings := sliceToString(d.Get("property_mappings_group").([]interface{}))
	setWrapper(d, "property_mappings_group", stringListConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	return diags
}

func resourceSourceLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapUpdate(ctx, d.Id()).LDAPSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
