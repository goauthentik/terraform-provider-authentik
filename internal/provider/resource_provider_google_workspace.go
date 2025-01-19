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

func resourceProviderGoogleWorkspace() *schema.Resource {
	return &schema.Resource{
		Description:   "Applications --- ",
		CreateContext: resourceProviderGoogleWorkspaceCreate,
		ReadContext:   resourceProviderGoogleWorkspaceRead,
		UpdateContext: resourceProviderGoogleWorkspaceUpdate,
		DeleteContext: resourceProviderGoogleWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      "JSON format expected. Use jsonencode() to pass objects.",
				DiffSuppressFunc: diffSuppressJSON,
			},
			"delegated_subject": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "seconds=0",
			},
			"default_group_email_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"property_mappings": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"property_mappings_group": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"exclude_users_service_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_delete_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description:      EnumToDescription(api.AllowedOutgoingSyncDeleteActionEnumValues),
				ValidateDiagFunc: StringInEnum(api.AllowedOutgoingSyncDeleteActionEnumValues),
			},
			"group_delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description: EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: StringInEnum([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
			},
		},
	}
}

func resourceProviderGoogleWorkspaceSchemaToProvider(d *schema.ResourceData) (*api.GoogleWorkspaceProviderRequest, diag.Diagnostics) {
	r := api.GoogleWorkspaceProviderRequest{
		Name:                       d.Get("name").(string),
		DelegatedSubject:           d.Get("delegated_subject").(string),
		DefaultGroupEmailDomain:    d.Get("default_group_email_domain").(string),
		PropertyMappings:           castSlice[string](d.Get("property_mappings").([]interface{})),
		PropertyMappingsGroup:      castSlice[string](d.Get("property_mappings_group").([]interface{})),
		ExcludeUsersServiceAccount: api.PtrBool(d.Get("exclude_users_service_account").(bool)),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(d.Get("user_delete_action").(string)).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(d.Get("group_delete_action").(string)).Ptr(),
	}

	if l, ok := d.Get("filter_group").(string); ok {
		r.FilterGroup = *api.NewNullableString(&l)
	}
	credentials := make(map[string]interface{})
	if l, ok := d.Get("credentials").(string); ok && l != "" {
		err := json.NewDecoder(strings.NewReader(l)).Decode(&credentials)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}
	r.Credentials = credentials
	return &r, nil
}

func resourceProviderGoogleWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderGoogleWorkspaceSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersGoogleWorkspaceCreate(ctx).GoogleWorkspaceProviderRequest(*r).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderGoogleWorkspaceRead(ctx, d, m)
}

func resourceProviderGoogleWorkspaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	res, hr, err := c.client.ProvidersApi.ProvidersGoogleWorkspaceRetrieve(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	setWrapper(d, "name", res.Name)
	setWrapper(d, "delegated_subject", res.DelegatedSubject)
	setWrapper(d, "default_group_email_domain", res.DefaultGroupEmailDomain)
	setWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	setWrapper(d, "user_delete_action", res.UserDeleteAction)
	setWrapper(d, "group_delete_action", res.GroupDeleteAction)
	setWrapper(d, "filter_group", res.FilterGroup)
	localMappings := castSlice[string](d.Get("property_mappings").([]interface{}))
	if len(localMappings) > 0 {
		setWrapper(d, "property_mappings", listConsistentMerge(localMappings, res.PropertyMappings))
	}
	localGroupMappings := castSlice[string](d.Get("property_mappings_group").([]interface{}))
	if len(localGroupMappings) > 0 {
		setWrapper(d, "property_mappings_group", listConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	}
	b, err := json.Marshal(res.Credentials)
	if err != nil {
		return diag.FromErr(err)
	}
	setWrapper(d, "credentials", string(b))
	return diags
}

func resourceProviderGoogleWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	app, diags := resourceProviderGoogleWorkspaceSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersGoogleWorkspaceUpdate(ctx, int32(id)).GoogleWorkspaceProviderRequest(*app).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}

	d.SetId(strconv.Itoa(int(res.Pk)))
	return resourceProviderGoogleWorkspaceRead(ctx, d, m)
}

func resourceProviderGoogleWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)
	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}
	hr, err := c.client.ProvidersApi.ProvidersGoogleWorkspaceDestroy(ctx, int32(id)).Execute()
	if err != nil {
		return httpToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
