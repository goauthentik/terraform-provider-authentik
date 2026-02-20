package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
			"sync_page_timeout": {
				Type:             schema.TypeString,
				Default:          "minutes=30",
				Optional:         true,
				ValidateDiagFunc: helpers.ValidateRelativeDuration,
				Description:      helpers.RelativeDurationDescription,
			},
			"sync_page_size": {
				Type:     schema.TypeInt,
				Default:  100,
				Optional: true,
			},
		},
	}
}

func resourceProviderGoogleWorkspaceSchemaToProvider(d *schema.ResourceData) (*api.GoogleWorkspaceProviderRequest, diag.Diagnostics) {
	r := api.GoogleWorkspaceProviderRequest{
		Name:                       d.Get("name").(string),
		DelegatedSubject:           d.Get("delegated_subject").(string),
		DefaultGroupEmailDomain:    d.Get("default_group_email_domain").(string),
		PropertyMappings:           helpers.CastSlice[string](d, "property_mappings"),
		PropertyMappingsGroup:      helpers.CastSlice[string](d, "property_mappings_group"),
		ExcludeUsersServiceAccount: new(d.Get("exclude_users_service_account").(bool)),
		UserDeleteAction:           api.OutgoingSyncDeleteAction(d.Get("user_delete_action").(string)).Ptr(),
		GroupDeleteAction:          api.OutgoingSyncDeleteAction(d.Get("group_delete_action").(string)).Ptr(),
		FilterGroup:                *api.NewNullableString(helpers.GetP[string](d, "filter_group")),
		DryRun:                     new(d.Get("dry_run").(bool)),
		SyncPageTimeout:            helpers.GetP[string](d, "sync_page_timeout"),
		SyncPageSize:               helpers.GetIntP(d, "sync_page_size"),
	}

	credentials, err := helpers.GetJSON[map[string]any](d, ("credentials"))
	r.Credentials = credentials
	return &r, err
}

func resourceProviderGoogleWorkspaceCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceProviderGoogleWorkspaceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	helpers.SetWrapper(d, "property_mappings", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings"),
		res.PropertyMappings,
	))
	helpers.SetWrapper(d, "property_mappings_group", helpers.ListConsistentMerge(
		helpers.CastSlice[string](d, "property_mappings_group"),
		res.PropertyMappingsGroup,
	))
	helpers.SetWrapper(d, "sync_page_timeout", res.SyncPageTimeout)
	helpers.SetWrapper(d, "sync_page_size", res.SyncPageSize)
	return helpers.SetJSON(d, "credentials", res.Credentials)
}

func resourceProviderGoogleWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceProviderGoogleWorkspaceDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
