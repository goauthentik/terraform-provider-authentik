package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
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
		Name:    d.Get("name").(string),
		Slug:    d.Get("slug").(string),
		Enabled: boolToPointer(d.Get("enabled").(bool)),

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

	propertyMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	r.PropertyMappings = &propertyMappings

	propertyMappingsGroup := sliceToString(d.Get("property_mappings_group").([]interface{}))
	r.PropertyMappingsGroup = &propertyMappingsGroup

	return &r
}

func resourceSourceLDAPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapCreate(ctx).LDAPSourceRequest(*r).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	res, hr, err := c.client.SourcesApi.SourcesLdapRetrieve(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	d.Set("uuid", res.Pk)
	d.Set("enabled", res.Enabled)

	d.Set("base_dn", res.BaseDn)
	d.Set("server_uri", res.ServerUri)
	d.Set("bind_cn", res.BindCn)
	d.Set("start_tls", res.StartTls)
	d.Set("base_dn", res.BaseDn)
	d.Set("additional_user_dn", res.AdditionalUserDn)
	d.Set("additional_group_dn", res.AdditionalGroupDn)
	d.Set("user_object_filter", res.UserObjectFilter)
	d.Set("group_object_filter", res.GroupObjectFilter)
	d.Set("group_membership_field", res.GroupMembershipField)
	d.Set("object_uniqueness_field", res.ObjectUniquenessField)
	d.Set("sync_users", res.SyncUsers)
	d.Set("sync_users_password", res.SyncUsersPassword)
	d.Set("sync_groups", res.SyncGroups)
	if res.SyncParentGroup.IsSet() {
		d.Set("sync_parent_group", res.SyncParentGroup.Get())
	}
	localMappings := sliceToString(d.Get("property_mappings").([]interface{}))
	d.Set("property_mappings", typeListConsistentMerge(localMappings, *res.PropertyMappings))
	localGroupMappings := sliceToString(d.Get("property_mappings_group").([]interface{}))
	d.Set("property_mappings_group", typeListConsistentMerge(localGroupMappings, *res.PropertyMappingsGroup))
	return diags
}

func resourceSourceLDAPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	app := resourceSourceLDAPSchemaToSource(d)

	res, hr, err := c.client.SourcesApi.SourcesLdapUpdate(ctx, d.Id()).LDAPSourceRequest(*app).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}

	d.SetId(res.Slug)
	return resourceSourceLDAPRead(ctx, d, m)
}

func resourceSourceLDAPDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	hr, err := c.client.SourcesApi.SourcesLdapDestroy(ctx, d.Id()).Execute()
	if err != nil {
		return httpToDiag(hr, err)
	}
	return diag.Diagnostics{}
}
