package provider

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/helpers"
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
			"dry_run": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"credentials": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				Description:      helpers.JSONDescription,
				DiffSuppressFunc: helpers.DiffSuppressJSON,
				ValidateDiagFunc: helpers.ValidateJSON,
			},
			"delegated_subject": {
				Type:     schema.TypeString,
				Optional: true,
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
				Description:      helpers.EnumToDescription(api.AllowedOutgoingSyncDeleteActionEnumValues),
				ValidateDiagFunc: helpers.StringInEnum(api.AllowedOutgoingSyncDeleteActionEnumValues),
			},
			"group_delete_action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  api.OUTGOINGSYNCDELETEACTION_DELETE,
				Description: helpers.EnumToDescription([]api.OutgoingSyncDeleteAction{
					api.OUTGOINGSYNCDELETEACTION_DELETE,
					api.OUTGOINGSYNCDELETEACTION_DO_NOTHING,
				}),
				ValidateDiagFunc: helpers.StringInEnum([]api.OutgoingSyncDeleteAction{
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
		PropertyMappings:           helpers.CastSlice[string](d.Get("property_mappings").([]interface{})),
		PropertyMappingsGroup:      helpers.CastSlice[string](d.Get("property_mappings_group").([]interface{})),
		ExcludeUsersServiceAccount: api.PtrBool(d.Get("exclude_users_service_account").(bool)),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(d.Get("user_delete_action").(string)).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(d.Get("group_delete_action").(string)).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.GetP[string](d, "filter_group")),
		DryRun:                     api.PtrBool(d.Get("dry_run").(bool)),
	}

	credentials, err := helpers.GetJSON[map[string]interface{}](d, ("credentials"))
	r.Credentials = credentials
	return &r, err
}

func resourceProviderGoogleWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*APIClient)

	r, diags := resourceProviderGoogleWorkspaceSchemaToProvider(d)
	if diags != nil {
		return diags
	}

	res, hr, err := c.client.ProvidersApi.ProvidersGoogleWorkspaceCreate(ctx).GoogleWorkspaceProviderRequest(*r).Execute()
	if err != nil {
		return helpers.HTTPToDiag(d, hr, err)
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
		return helpers.HTTPToDiag(d, hr, err)
	}

	helpers.SetWrapper(d, "name", res.Name)
	helpers.SetWrapper(d, "delegated_subject", res.DelegatedSubject)
	helpers.SetWrapper(d, "default_group_email_domain", res.DefaultGroupEmailDomain)
	helpers.SetWrapper(d, "exclude_users_service_account", res.ExcludeUsersServiceAccount)
	helpers.SetWrapper(d, "user_delete_action", res.UserDeleteAction)
	helpers.SetWrapper(d, "group_delete_action", res.GroupDeleteAction)
	helpers.SetWrapper(d, "filter_group", res.FilterGroup)
	helpers.SetWrapper(d, "dry_run", res.DryRun)
	localMappings := helpers.CastSlice[string](d.Get("property_mappings").([]interface{}))
	if len(localMappings) > 0 {
		helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(localMappings, res.PropertyMappings))
	}
	localGroupMappings := helpers.CastSlice[string](d.Get("property_mappings_group").([]interface{}))
	if len(localGroupMappings) > 0 {
		helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(localGroupMappings, res.PropertyMappingsGroup))
	}
	b, err := json.Marshal(res.Credentials)
	if err != nil {
		return diag.FromErr(err)
	}
	helpers.SetWrapper(d, "credentials", string(b))
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
		return helpers.HTTPToDiag(d, hr, err)
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
		return helpers.HTTPToDiag(d, hr, err)
	}
	return diag.Diagnostics{}
}
